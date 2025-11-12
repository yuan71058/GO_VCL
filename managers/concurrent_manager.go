// managers/concurrent_manager.go - 多线程并发管理器
// 该文件实现了多线程并发任务管理器，提供工作者池、任务调度、结果收集和性能监控功能
// 支持多种类型的并发任务执行，包括HTTP请求、计算密集型任务和I/O操作
package managers

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"windows-gui-app/interfaces"
)

// TaskResult 任务执行结果
// 封装了任务执行后的详细信息，包括任务ID、名称、成功状态、结果数据和错误信息
type TaskResult struct {
	ID       int         // 任务唯一标识符
	TaskName string      // 任务名称，用于显示和识别
	Success  bool        // 任务执行是否成功
	Result   interface{} // 任务执行结果数据
	Error    error       // 任务执行错误信息（如果有）
	Duration time.Duration // 任务执行耗时
}

// ConcurrentManager 并发管理器
// 提供多线程并发任务管理功能，包括工作者池、任务调度、结果收集和性能监控
// 使用通道和互斥锁确保线程安全，支持UI线程安全操作
type ConcurrentManager struct {
	ui            interfaces.UIInterface // UI实例接口，用于更新界面状态和数据
	workerPool    chan struct{}         // 工作者池通道，控制并发任务数量
	results       []TaskResult          // 任务执行结果集合
	resultsMutex  sync.RWMutex          // 保护结果集合的读写锁
	taskCounter   int                   // 任务计数器，用于生成唯一任务ID
	counterMutex  sync.Mutex            // 保护任务计数器的互斥锁
	workerCount   int                   // 工作线程数量
	uiChannel     chan func()           // UI操作同步通道，确保UI操作在主线程执行
	uiClose       chan bool             // UI工作协程关闭通道
}

// NewConcurrentManager 创建新的并发管理器
// 初始化并发管理器，创建工作者池和UI同步通道，启动UI工作协程
// 工作线程数量设置为CPU核心数的2倍，以充分利用CPU资源
// 参数:
//   - ui: UI接口实例，用于更新界面状态（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *ConcurrentManager: 并发管理器实例
func NewConcurrentManager(ui interfaces.UIInterface) *ConcurrentManager {
	workerCount := runtime.NumCPU() * 2 // 使用CPU核心数 * 2作为工作线程数
	cm := &ConcurrentManager{
		ui:            ui,
		workerPool:    make(chan struct{}, workerCount),
		results:       make([]TaskResult, 0),
		workerCount:   workerCount,
		taskCounter:   0,
		uiChannel:     make(chan func(), 1000), // 缓冲区大小为1000，防止阻塞
		uiClose:       make(chan bool),
	}
	
	// 启动UI操作协程来处理UI调用，确保UI操作在主线程执行
	go cm.uiWorker()
	
	return cm
}

// SetUIInstance 设置UI实例
// 允许在创建管理器后设置UI实例，用于解耦合UI和管理器的创建顺序
// 参数:
//   - ui: UI实例接口
func (cm *ConcurrentManager) SetUIInstance(ui interfaces.UIInterface) {
	cm.ui = ui
}

// RunConcurrentTasks 执行并发任务示例
// 根据任务类型执行不同类型的并发任务，支持HTTP请求、计算密集型任务和I/O操作
// 记录任务执行开始和结束时间，统计成功和失败任务数量
// 参数:
//   - taskType: 任务类型，可选值为"http"、"compute"、"io"或"mixed"
//   - count: 任务数量，要执行的任务总数
//
// 返回:
//   - error: 操作错误，nil表示成功
func (cm *ConcurrentManager) RunConcurrentTasks(taskType string, count int) error {
	cm.safeAddLog(fmt.Sprintf("=== 开始执行 %d 个并发任务 (类型: %s) ===", count, taskType))

	startTime := time.Now()

	// 根据任务类型调用相应的处理函数
	switch taskType {
	case "http":
		return cm.runConcurrentHTTPTasks(count)
	case "compute":
		return cm.runConcurrentComputeTasks(count)
	case "io":
		return cm.runConcurrentIOTasks(count)
	case "mixed":
		return cm.runConcurrentMixedTasks(count)
	default:
		return fmt.Errorf("不支持的任务类型: %s", taskType)
	}

	duration := time.Since(startTime)
	cm.safeAddLog(fmt.Sprintf("所有任务完成，总耗时: %v", duration))
	cm.safeAddLog(fmt.Sprintf("成功任务: %d, 失败任务: %d", cm.getSuccessCount(), cm.getFailureCount()))

	return nil
}

// runConcurrentHTTPTasks 执行并发HTTP请求任务
// 创建指定数量的并发HTTP请求任务，模拟网络请求操作
// 使用工作者池控制并发数量，每个任务随机选择一个URL进行请求
// 80%的任务会成功，20%的任务会失败，模拟真实网络环境
// 参数:
//   - count: 要执行的HTTP任务数量
//
// 返回:
//   - error: 操作错误，nil表示成功
func (cm *ConcurrentManager) runConcurrentHTTPTasks(count int) error {
	var wg sync.WaitGroup
	// 预定义的URL列表，循环使用
	urls := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
		"https://jsonplaceholder.typicode.com/users/1",
		"https://jsonplaceholder.typicode.com/users/2",
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		// 获取一个工作线程槽位，如果池满则阻塞等待
		cm.workerPool <- struct{}{}
		
		taskID := cm.getNextTaskID()
		taskName := fmt.Sprintf("HTTP请求_%d", taskID)
		url := urls[i%len(urls)]
		
		// 启动goroutine执行任务
		go func(taskID int, taskName, url string) {
			defer wg.Done()
			defer func() { <-cm.workerPool }() // 释放工作线程槽位

			startTime := time.Now()
			result := TaskResult{
				ID:       taskID,
				TaskName: taskName,
				Success:  false,
			}

			// 确保任务结果被记录
			defer func() {
				result.Duration = time.Since(startTime)
				cm.addResult(result)
			}()

			cm.safeAddLog(fmt.Sprintf("[%s] 开始执行...", taskName))

			// 模拟HTTP请求，随机延迟100-600毫秒
			time.Sleep(time.Duration(100+rand.Intn(500)) * time.Millisecond)
			
			// 80%成功率，20%失败率，模拟真实网络环境
			if rand.Intn(10) < 8 {
				result.Success = true
				result.Result = fmt.Sprintf("HTTP响应数据 (长度: %d)", rand.Intn(1000)+100)
				cm.safeAddLog(fmt.Sprintf("[%s] 成功完成", taskName))
			} else {
				result.Error = fmt.Errorf("网络超时")
				cm.safeAddLog(fmt.Sprintf("[%s] 失败: %v", taskName, result.Error))
			}
		}(taskID, taskName, url)
	}

	wg.Wait()
	return nil
}

// runConcurrentComputeTasks 执行并发计算任务
// 创建指定数量的并发计算密集型任务，模拟CPU密集型工作负载
// 每个任务计算斐波那契数列，工作负载随机，定期让出CPU时间片
// 参数:
//   - count: 要执行的计算任务数量
//
// 返回:
//   - error: 操作错误，nil表示成功
func (cm *ConcurrentManager) runConcurrentComputeTasks(count int) error {
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		cm.workerPool <- struct{}{}
		
		taskID := cm.getNextTaskID()
		taskName := fmt.Sprintf("计算任务_%d", taskID)
		workload := rand.Intn(1000) + 100 // 随机工作负载100-1100
		
		go func(taskID int, taskName string, workload int) {
			defer wg.Done()
			defer func() { <-cm.workerPool }()

			startTime := time.Now()
			result := TaskResult{
				ID:       taskID,
				TaskName: taskName,
				Success:  false,
			}

			defer func() {
				result.Duration = time.Since(startTime)
				cm.addResult(result)
			}()

			cm.safeAddLog(fmt.Sprintf("[%s] 开始执行 (工作量: %d)", taskName, workload))

			// 模拟计算任务 - 斐波那契数列计算
			sum := 0
			for j := 0; j < workload; j++ {
				sum += fibonacci(j % 20) // 限制计算复杂度，防止计算时间过长
				// 偶尔让出CPU时间片，避免长时间占用CPU
				if j%100 == 0 {
					runtime.Gosched()
				}
			}

			result.Success = true
			result.Result = sum
			
			cm.safeAddLog(fmt.Sprintf("[%s] 计算完成，结果: %d", taskName, sum))
		}(taskID, taskName, workload)
	}

	wg.Wait()
	return nil
}

// runConcurrentIOTasks 执行并发I/O任务
// 创建指定数量的并发I/O任务，模拟文件读写操作
// 每个任务处理随机大小的数据块，模拟I/O密集型工作负载
// 参数:
//   - count: 要执行的I/O任务数量
//
// 返回:
//   - error: 操作错误，nil表示成功
func (cm *ConcurrentManager) runConcurrentIOTasks(count int) error {
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		cm.workerPool <- struct{}{}
		
		taskID := cm.getNextTaskID()
		taskName := fmt.Sprintf("I/O任务_%d", taskID)
		fileSize := rand.Intn(1024*1024) + 1024 // 1KB到1MB的随机文件大小
		
		go func(taskID int, taskName string, fileSize int) {
			defer wg.Done()
			defer func() { <-cm.workerPool }()

			startTime := time.Now()
			result := TaskResult{
				ID:       taskID,
				TaskName: taskName,
				Success:  false,
			}

			defer func() {
				result.Duration = time.Since(startTime)
				cm.addResult(result)
			}()

			cm.safeAddLog(fmt.Sprintf("[%s] 开始执行 (文件大小: %d bytes)", taskName, fileSize))

			// 模拟文件I/O操作，创建随机数据
			data := make([]byte, fileSize)
			for i := 0; i < fileSize; i++ {
				data[i] = byte(rand.Intn(256))
				// 模拟读写时间，每10000字节暂停1微秒
				if i%10000 == 0 {
					time.Sleep(1 * time.Microsecond)
				}
			}

			result.Success = true
			result.Result = fmt.Sprintf("处理了 %d bytes 数据", fileSize)
			
			cm.safeAddLog(fmt.Sprintf("[%s] I/O操作完成", taskName))
		}(taskID, taskName, fileSize)
	}

	wg.Wait()
	return nil
}

// runConcurrentMixedTasks 执行混合类型并发任务
// 创建指定数量的混合类型并发任务，包括HTTP请求、计算任务和I/O任务
// 任务类型循环分配，模拟真实应用中的混合工作负载
// 参数:
//   - count: 要执行的混合任务数量
//
// 返回:
//   - error: 操作错误，nil表示成功
func (cm *ConcurrentManager) runConcurrentMixedTasks(count int) error {
	var wg sync.WaitGroup
	// 任务类型列表，循环使用
	taskTypes := []string{"http", "compute", "io"}

	for i := 0; i < count; i++ {
		wg.Add(1)
		cm.workerPool <- struct{}{}
		
		taskID := cm.getNextTaskID()
		taskType := taskTypes[i%len(taskTypes)]
		taskName := fmt.Sprintf("混合任务_%d (%s)", taskID, taskType)
		
		go func(taskID int, taskType, taskName string) {
			defer wg.Done()
			defer func() { <-cm.workerPool }()

			startTime := time.Now()
			result := TaskResult{
				ID:       taskID,
				TaskName: taskName,
				Success:  false,
			}

			defer func() {
				result.Duration = time.Since(startTime)
				cm.addResult(result)
			}()

			cm.safeAddLog(fmt.Sprintf("[%s] 开始执行 (类型: %s)", taskName, taskType))

			// 根据任务类型执行不同操作
			switch taskType {
			case "http":
				// 模拟HTTP请求
				time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
				if rand.Intn(10) < 8 {
					result.Success = true
					result.Result = "HTTP响应数据"
				} else {
					result.Error = fmt.Errorf("请求失败")
				}
			case "compute":
				// 模拟计算任务
				sum := fibonacci(rand.Intn(20))
				result.Success = true
				result.Result = sum
			case "io":
				// 模拟I/O操作
				time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
				result.Success = true
				result.Result = "文件处理完成"
			}

			if result.Success {
				cm.safeAddLog(fmt.Sprintf("[%s] 成功完成", taskName))
			} else {
				cm.safeAddLog(fmt.Sprintf("[%s] 失败: %v", taskName, result.Error))
			}
		}(taskID, taskType, taskName)
	}

	wg.Wait()
	return nil
}

// runWorkerPoolTask 使用工作者池执行任务
// 尝试获取工作者池槽位并执行任务，如果池满则立即返回错误
// 参数:
//   - taskName: 任务名称
//   - taskFunc: 要执行的任务函数
//
// 返回:
//   - error: 操作错误，nil表示成功
func (cm *ConcurrentManager) runWorkerPoolTask(taskName string, taskFunc func() error) error {
	select {
	case cm.workerPool <- struct{}{}: // 获取工作线程槽位
		defer func() { <-cm.workerPool }()
		return taskFunc()
	default:
		return fmt.Errorf("工作者池已满，请稍后再试")
	}
}

// ExecuteTaskWithContext 使用上下文执行可取消的任务
// 支持任务取消功能，当上下文被取消时，任务可以提前终止
// 记录任务执行时间和结果，支持通过上下文控制任务生命周期
// 参数:
//   - ctx: 上下文，用于控制任务取消和超时
//   - taskName: 任务名称
//   - taskFunc: 要执行的任务函数，接收上下文参数
//
// 返回:
//   - error: 任务执行错误，nil表示成功
func (cm *ConcurrentManager) ExecuteTaskWithContext(ctx context.Context, taskName string, taskFunc func(context.Context) error) error {
	result := TaskResult{
		ID:       cm.getNextTaskID(),
		TaskName: taskName,
		Success:  false,
	}

	startTime := time.Now()

	cm.safeAddLog(fmt.Sprintf("[%s] 开始执行 (可取消)", taskName))

	err := taskFunc(ctx)

	result.Duration = time.Since(startTime)
	if err != nil {
		result.Error = err
		cm.safeAddLog(fmt.Sprintf("[%s] 失败: %v", taskName, err))
	} else {
		result.Success = true
		cm.safeAddLog(fmt.Sprintf("[%s] 成功完成", taskName))
	}

	cm.addResult(result)
	return err
}

// GetResults 获取所有任务结果
// 返回所有已执行任务的详细结果，包括成功和失败的任务
// 使用读锁确保线程安全，返回结果副本防止外部修改
// 返回:
//   - []TaskResult: 任务结果列表
func (cm *ConcurrentManager) GetResults() []TaskResult {
	cm.resultsMutex.RLock()
	defer cm.resultsMutex.RUnlock()
	return append([]TaskResult(nil), cm.results...)
}

// GetStatistics 获取任务统计信息
// 计算并返回任务执行的统计信息，包括总数、成功数、失败数、成功率等
// 使用读锁确保线程安全，计算平均执行时间等性能指标
// 返回:
//   - map[string]interface{}: 包含统计信息的字典
func (cm *ConcurrentManager) GetStatistics() map[string]interface{} {
	cm.resultsMutex.RLock()
	defer cm.resultsMutex.RUnlock()

	total := len(cm.results)
	success := 0
	failure := 0
	var totalDuration time.Duration

	// 遍历所有结果，统计成功和失败任务数量
	for _, result := range cm.results {
		if result.Success {
			success++
		} else {
			failure++
		}
		totalDuration += result.Duration
	}

	var avgDuration time.Duration
	if total > 0 {
		avgDuration = totalDuration / time.Duration(total)
	}

	return map[string]interface{}{
		"total_tasks":      total,
		"success_tasks":    success,
		"failure_tasks":    failure,
		"success_rate":     float64(success) / float64(total) * 100,
		"total_duration":   totalDuration,
		"average_duration": avgDuration,
		"worker_count":     cm.workerCount,
	}
}

// GetPerformanceMetrics 获取性能指标
// 获取系统性能指标，包括吞吐量、CPU数量、goroutine数量和内存使用情况
// 这些指标可用于评估并发任务执行的性能和资源利用率
// 返回:
//   - map[string]interface{}: 包含性能指标的字典
func (cm *ConcurrentManager) GetPerformanceMetrics() map[string]interface{} {
	stats := cm.GetStatistics()
	
	// 计算吞吐量（每秒完成的任务数）
	throughput := 0.0
	if stats["total_duration"].(time.Duration).Seconds() > 0 {
		throughput = float64(stats["success_tasks"].(int)) / stats["total_duration"].(time.Duration).Seconds()
	}

	// 获取内存统计信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"throughput_per_second": throughput,
		"cpu_count":             runtime.NumCPU(),
		"goroutine_count":       runtime.NumGoroutine(),
		"memory_alloc":          m.Alloc / 1024 / 1024, // MB
		"worker_pool_size":      cm.workerCount,
	}
}

// Reset 重置管理器状态
// 清空所有任务结果和计数器，重置管理器到初始状态
// 使用写锁确保线程安全，在重置过程中阻止其他操作
func (cm *ConcurrentManager) Reset() {
	cm.resultsMutex.Lock()
	defer cm.resultsMutex.Unlock()
	cm.results = make([]TaskResult, 0)
	cm.taskCounter = 0
}

// 辅助函数

// getNextTaskID 获取下一个任务ID
// 线程安全地生成递增的任务ID，确保每个任务有唯一标识
// 返回:
//   - int: 新的任务ID
func (cm *ConcurrentManager) getNextTaskID() int {
	cm.counterMutex.Lock()
	defer cm.counterMutex.Unlock()
	cm.taskCounter++
	return cm.taskCounter
}

// addResult 添加任务结果
// 线程安全地将任务结果添加到结果集合中
// 参数:
//   - result: 要添加的任务结果
func (cm *ConcurrentManager) addResult(result TaskResult) {
	cm.resultsMutex.Lock()
	defer cm.resultsMutex.Unlock()
	cm.results = append(cm.results, result)
}

// getSuccessCount 获取成功任务数量
// 线程安全地统计成功完成的任务数量
// 返回:
//   - int: 成功任务数量
func (cm *ConcurrentManager) getSuccessCount() int {
	cm.resultsMutex.RLock()
	defer cm.resultsMutex.RUnlock()
	
	count := 0
	for _, result := range cm.results {
		if result.Success {
			count++
		}
	}
	return count
}

// getFailureCount 获取失败任务数量
// 线程安全地统计失败的任务数量
// 返回:
//   - int: 失败任务数量
func (cm *ConcurrentManager) getFailureCount() int {
	cm.resultsMutex.RLock()
	defer cm.resultsMutex.RUnlock()
	
	count := 0
	for _, result := range cm.results {
		if !result.Success {
			count++
		}
	}
	return count
}

// uiWorker UI操作专用工作协程
// 持续监听UI操作通道，在主线程中执行UI更新操作
// 确保所有UI操作都在同一个线程中执行，避免并发UI操作导致的问题
func (cm *ConcurrentManager) uiWorker() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("UI工作协程异常: %v\n", r)
		}
	}()

	for {
		select {
		case task := <-cm.uiChannel:
			if task != nil {
				task()
			}
		case <-cm.uiClose:
			return
		}
	}
}

// safeAddLog 线程安全的日志添加方法
// 通过UI操作通道安全地添加日志，确保UI操作在主线程执行
// 如果通道缓冲区已满，则将日志输出到控制台，避免阻塞
// 参数:
//   - logText: 要添加的日志文本
func (cm *ConcurrentManager) safeAddLog(logText string) {
	if cm.ui == nil {
		return
	}
	
	select {
	case cm.uiChannel <- func() {
		cm.ui.AddLog(logText)
	}:
	default:
		fmt.Printf("UI通道缓冲区已满，无法添加日志: %s\n", logText)
	}
}

// Close 关闭并发管理器
// 关闭UI工作协程和相关通道，释放资源
// 应该在不再使用管理器时调用，防止资源泄漏
func (cm *ConcurrentManager) Close() {
	close(cm.uiClose)
	close(cm.uiChannel)
}

// fibonacci 斐波那契数列计算（递归）
// 计算第n个斐波那契数，用于模拟计算密集型任务
// 注意：递归实现效率较低，仅用于演示目的
// 参数:
//   - n: 要计算的斐波那契数列索引
//
// 返回:
//   - int: 第n个斐波那契数
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}