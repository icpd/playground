> taskset 是一个 Linux 命令行工具,用于设置或检索一个进程的 CPU 亲和性(CPU affinity)。CPU 亲和性决定了一个进程可以运行在哪个 CPU 或 CPU 核心上。 合理设置 CPU 亲和性可以带来以下好处:
>
> 减少缓存丢失,提高缓存命中率  
> 降低线程迁移和上下文切换开销  
> 分配独占资源,避免资源争用  
> 均衡负载,充分利用多核 CPU  


# taskset 使用方法
## 获取进程的 CPU 亲和性
`taskset -p <PID>`  
该命令会显示进程 ID 为 <PID> 的进程当前的 CPU 亲和性掩码 (affinity mask)。

## 设置进程的 CPU 亲和性
`taskset -p <MASK> <PID>`  
该命令将进程 ID 为 <PID> 的进程的 CPU 亲和性设置为 <MASK> 指定的 CPU 掩码。

CPU 掩码是一个十六进制数, 用于表示要被绑定的 CPU 核心。例如, 在一个 4 核心的系统中, 0x1 表示只绑定到第一个 CPU 核心, 0x5 表示绑定到第 1 个和第 3 个 CPU 核心。哪一个位上为 1，就会绑定到哪个核上。

## 启动新进程并设置其 CPU 亲和性
`taskset <MASK> <COMMAND> [ARGUMENTS...]`  
该命令会使用 <MASK> 指定的 CPU 亲和性启动一个新进程, 运行命令 <COMMAND> 及其参数 [ARGUMENTS...]。

## 其他选项
taskset 还提供了其他一些有用的选项:

- c: 使用 CPU 列表格式指定 CPU 掩码, 例如 taskset -c 0,5,7,9-11 ./program  
- a: 获取或设置进程及其所有线程的 CPU 亲和性  
例子：
```
默认行为是运行一条新命令：  
taskset 03 sshd -b 1024  
您可以获取现有任务的掩码：  
taskset -p 700   
或设置掩码：    
taskset -p 03 700  
使用逗号分隔的列表格式而不是掩码：  
taskset -pc 0,3,7-11 700  
列表格式中的范围可以带一个跨度参数：  
例如 0-31:2 与掩码 0x55555555 等效  
```
## 子进程的 CPU 亲和性  
Linux 默认子进程会继承父进程的 CPU 亲和性,以保持亲和性的连续性。

但如果有需要, 用户依然可以手动设置属于它的子进程或线程的 CPU 亲和性,以优化 CPU 资源利用和程序性能。

如果想修改其他用户的进程的 CPU 亲和性，需要 root 权限或者拥有 CAP_SYS_NICE 权限。

任何用户都可以获取任意进程的 CPU 亲和性掩码。