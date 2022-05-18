# Socket
![socket_conn](img/socket_conn.png)

# Server Side

- Listen Function
    
    ```fsharp
    func Listen(network, address string) (Listener, error)
    
    network: protocol like TCP, UDP
    address: IP Address + Port Number
    ```
    
- Listener Interface
    
    ```fsharp
    type Listener interface {
    		Accept() (Conn, error)
    		Close() error
    		Addr() Addr
    }
    ```
    
- Conn Interface
    
    ```fsharp
    type Conn interface {
    		Read(b []byte) (n int, err error)
    		Write(b []byte) (n, int, err error)
    		Close() error
    		LocalAddr() Addr
    		RemoteAddr() Addr
    		SetDeadline(t time.Time) error
    		SetReadDeadline(t time.Time) error
    		SetWriteDeadline(t time.Time) error
    }
    ```
    

# TCP-CS
![tcp_cs](img/tcp_cs.png)

## Server

```go
func main(){
		// define server protocal and port number
		listener, err := net.Listen("tcp","172.0.0.1:8000")
		if err != nil {
				fmt.Println("net.Listen err: ", err)
				return
		}

		defer listener.Close()

		fmt.Println("waiting client to build connection...")

		// blocking connetion requests from client
		// create connetion successfully and return socket to communication
		conn, err := listener.Accept()
		if err != nil {
				fmt.Println("listener.Accept() err: ", err)
				return
		}

		defer conn.Close()

		fmt.Println("server and client successfully build connection")

		// read client data
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
				fmt.Println("conn.Read() err: ", err)
				return
		}

		conn.Write(buf[:n])

		fmt.Println("server read client data: ",string(buf[:n]))

}
```

---

## Client
1. conn, err := net.Dial(”tcp”, IP + port)
2. write data to server → conn.Write()
3. Read server response data → conn.Read()
4. conn.Close()

```go
func main(){
		// define server protocal and port number
		conn, err := net.Dial("tcp","172.0.0.1:8000")
		if err != nil {
				fmt.Println("net.Dial err: ", err)
				return
		}

		defer conn.Close()

		conn.Write([]byte("Are you Ready?"))

		buf := make([]byte,4096)
		n, err := conn.Read(buf)
		if err != nil {
				fmt.Println("conn.Read() err: ", err)
				return
		}

		fmt.Println("server response data: ", string(buf[:n]))

}
```

---

## Server - Concurrent
1. create listener → listener := net.Listen(”tcp”, IP + port)
2. defer listener.Close()
3. for loop blocking to listen client connection → conn := listener.Accept()
4. create goroutine for each client request and transfer data → go HandlerConnect()
5. implement HandlerConnet(conn net.Conn)
    1. defer conn.Close()
    2. get client address → conn.RemoteAddr()
    3. handle data to upper → strings.ToUpper()
    4. response data → conn.Write(buf[:n])
6. server check if connection closed: **Read client data and return 0**

```go
func HandlerConnect(conn net.Conn){
		defer conn.Close()

		// get remote client address
		addr := conn.RemoteAddr()
		fmt.Println(addr,"build connection successfully")

		// loop read client data
		buf := make([]byte,4096)
		for {
				n,err := conn.Read(buf)
				if n == 0 || "exit\n" == string(buf[:n]){
						fmt.Println("client connection closed!")
						return
				}
				if err != nil {
						fmt.Println("conn.Read() err: ", err)
						return
				}

				// use data
				fmt.Println("Server read data successfully: ",string(buf[:n]))

				// toupper and response
				conn.Write([]byte(strings.ToUpper(string(buf[:n]))))
		}
}

func main(){
		// define server protocal and port number
		listener, err := net.Listen("tcp","172.0.0.1:8000")
		if err != nil {
				fmt.Println("net.Listen err: ", err)
				return
		}

		defer listener.Close()

		// blocking connetion requests from client
		// create connetion successfully and return socket to communication
		for {
				fmt.Println("waiting client to build connection...")

				conn, err := listener.Accept()
				if err != nil {
						fmt.Println("listener.Accept() err: ", err)
						return
				}
				// 具體完成 server 及 client 資料通訊請求
				go HandlerConnect(conn)
		}

		fmt.Println("server and client successfully build connection")

		// read client data
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
				fmt.Println("conn.Read() err: ", err)
				return
		}

		conn.Write(buf[:n])

		fmt.Println("server read client data: ",string(buf[:n]))

}
```

---

## Client - Concurrent

```go
func main(){
		// define server protocal and port number
		conn, err := net.Dial("tcp","172.0.0.1:8000")
		if err != nil {
				fmt.Println("net.Dial err: ", err)
				return
		}

		defer conn.Close()

		// 獲取 user 鍵盤輸入(stdin)，將輸入資料送給 server
		go func() {
				str := make([]byte, 4096)
				for {
						n, err := os.Stdin.Read(str)
						if err != nil {
								fmt.Println("os.Stdin.Read err: ", err)
								continue
						}
						// 寫給 server
						conn.Write(str[:n])
				}
		}()

		// 顯示 server response data
		buf := make([]byte, 4096)
		for {
				n, err := conn.Read(buf)
				if err != nil {
						fmt.Println("conn.Read err: ", err)
						return
				}
				fmt.Println("client read server response: ", string(buf[:n]))
		}

		buf := make([]byte,4096)
		n, err := conn.Read(buf)
		if n == 0 {
				fmt.Println("server connection closed!")
				return
		}
		if err != nil {
				fmt.Println("conn.Read() err: ", err)
				return
		}

		fmt.Println("server response data: ", string(buf[:n]))

}
```