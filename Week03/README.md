### 学习笔记

- errorgroup  全部出现完成  其中一个error 才会 wait 结束
所有实现优雅关闭 在监听 linux 信号的线程中 要调用 总的context 进行cancel
 
- server 启动函数
    ```
        func runServer(ctx context.Context, serverPort int, serverName string) error {
            s := &http.Server{
                Addr:           fmt.Sprintf(":%d", serverPort),
                MaxHeaderBytes: 1 << 20,
            }

            defer func() {
                if err := recover(); err != nil {
                    log.Printf("%+v",err)
                    log.Println(serverName,"into defer recover()")
                    cancelServer(s,serverName)
                }
            }()
            go func() {
                <-ctx.Done()
                cancelServer(s,serverName)
            }()
            log.Println(serverPort,serverName," with run")
            err := s.ListenAndServe()
            log.Println(serverName,"with error",err)
            return err
        }
    ```

- listen sign 函数
    ```
        func listenStopSignal(ctx context.Context, sig chan os.Signal, cancel context.CancelFunc) error {
            select {
            case sign := <- sig:
                log.Println("接受信号",sign)
                err := fmt.Errorf("handle signal: %d", sign)
                cancel()
                return err
            case <-ctx.Done():
                log.Println("listenStopSignal ctx.Done")
                return nil
            }

        }
    ```