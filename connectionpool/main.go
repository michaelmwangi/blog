package main

import (
	"fmt"
	"net"
	"time"
)

type TCPPool struct{
	pool []*net.TCPConn //tcp connections storage 

}

// put a tcp connection in the pool
func (t *TCPPool)putConnection(conn *net.TCPConn){
	t.pool = append(t.pool, conn)
}

// get and return a free tcp connection from the pool, nil if there are no free connections
func (t *TCPPool)getConnection() *net.TCPConn{
	num := len(t.pool)
	if num == 0{
		return nil
	}
	conn := t.pool[num-1]

	//pop last item from the pool 
	t.pool = t.pool[:num-1]
	return conn
}

func (t *TCPPool)closeConnections(){
	for _, con := range t.pool{
		con.Close()
	}
}

func createConnection(addr string) (*net.TCPConn, error){

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr);
	if err != nil{
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr); 
	if err != nil{
		return nil, err
	}
	return conn, nil

}

// initialize and return the tcppool with num connections
func newPool(addr string, num int) *TCPPool{
	if(num < 1){
		fmt.Println("Connection pool requires atleast one connection")
		return nil
	}
	tcpPool := &TCPPool{}
	for i := 0; i < num; i++{
		conn, err := createConnection(addr);
		if err != nil{
			fmt.Println("Cannot create connection ", err)
			return nil
		}
		tcpPool.putConnection(conn)	
	}
	return tcpPool
}

func send(conn *net.TCPConn, data []byte) (int, error){

	num, err := conn.Write(data);
	if err != nil{
		return 0, err
	}
	
	return num, nil
}

func read(conn *net.TCPConn, bufSize int)([]byte, error){
	buf := make([]byte, bufSize)
	if 	_, err := conn.Read(buf); err != nil{
		return nil, err
	}
	return buf, nil
}

func withPool(pool *TCPPool){

	start := time.Now()
	conn := pool.getConnection()
	if conn == nil{
		fmt.Println("Could not get a free connection")
		return 
	}

	defer pool.putConnection(conn)
	if _, err := send(conn, []byte("ping\r\n")); err != nil{
		fmt.Println(err)
		return
	}

	_, err := read(conn, 32)
	if err != nil{
		fmt.Println(err)
		return 
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println("pool ",diff)
}


func noPool(){
	startTime := time.Now()
	conn, err := createConnection("127.0.0.1:6379")
	if err != nil{
		fmt.Println(err)
	}
	defer conn.Close()
	if _, err := send(conn, []byte("ping\r\n")); err != nil{
		fmt.Println(err)
		return
	}

	_, err = read(conn, 32)
	if err != nil{
		fmt.Println(err)
		return 
	}

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("no pool ",diff)
}

func main() {
	fmt.Println("Starting")
	//use a connection pool
	mrPool := newPool("127.0.0.1:6379", 1)
	defer mrPool.closeConnections()
	withPool(mrPool)

	//create a connection everytime
	noPool()
}

