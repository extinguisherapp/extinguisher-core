package main

import (
	"bufio"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net"
	"net/http/httputil"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

type Project struct {
	Path           string
	HTMLSocketPath string
	Default        bool

	htmlCmd     *exec.Cmd
	serverTimer *time.Timer
	lock        sync.Mutex
}

func NewProject(path string) *Project {
	return &Project{Path: path, lock: sync.Mutex{}}
}

func (p *Project) Name() string {
	return path.Base(p.Path)
}
func (p *Project) StartService() {
	serverTTL := 180 * time.Second

	p.lock.Lock()
	if p.htmlCmd == nil {
		p.serverTimer = time.NewTimer(serverTTL)
		p.htmlCmd = exec.Command("ruby", "ruby/run.rb", p.Path)
		stdout, err := p.htmlCmd.StdoutPipe()
		checkErr(err)
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				log.Println(scanner.Text()) // Println will add back the final '\n'
			}
			if err := scanner.Err(); err != nil {
				log.Fatalln("reading standard output:", err)
			}
		}()

		stderr, err := p.htmlCmd.StderrPipe()
		checkErr(err)
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				log.Println(scanner.Text()) // Println will add back the final '\n'
			}
			if err := scanner.Err(); err != nil {
				log.Fatalln("reading standard error:", err)
			}
		}()

		err = p.htmlCmd.Start()
		checkErr(err)

		if err == nil {
			p.HTMLSocketPath = path.Join(p.Path, "html.socket")
		}

		go func() {
			<-p.serverTimer.C
			log.Println("stop " + p.Name() + " server")
			p.htmlCmd.Process.Kill()
			p.htmlCmd.Wait()
			p.htmlCmd = nil
		}()
	}
	p.serverTimer.Reset(serverTTL)
	p.lock.Unlock()
}

var (
	projects           []*Project
	ErrProjectNotFound = errors.New("Project Not Found")
)

func main() {
	projects = append(projects, NewProject("/home/tka/Project/Compass/fireapp-beginner-example"))
	projects = append(projects, NewProject("/home/tka/Project/Compass/fire-app-sample-project"))
	r := gin.Default()
	r.GET("/*any", requestHandler)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

func getProjectByName(name string) (*Project, error) {
	for _, p := range projects {
		if p.Name() == name {
			return p, nil
		}
	}
	return nil, ErrProjectNotFound
}

func requestHandler(c *gin.Context) {
	projectName := strings.Split(c.Request.Host, ".")[0]

	project, err := getProjectByName(projectName)
	if err == ErrProjectNotFound {
		c.String(404, "Not Found")
		return
	}
	checkErr(err)

	project.StartService()

	dump, err := httputil.DumpRequest(c.Request, false)
	checkErr(err)

	/*
	* try to connect html socket server
	* if fail 6 times,  show error message
	 */
	var maxRetry int
	var sock net.Conn
	for maxRetry = 10; maxRetry > 0; maxRetry-- {
		sock, err = net.Dial("unix", project.HTMLSocketPath)
		if err != nil {
			time.Sleep(250 * time.Millisecond)
			checkErr(err)
		} else {
			break
		}
	}
	if maxRetry == 0 {
		c.String(500, "Cant Connect To HTML Server\n"+err.Error())
		return
	}

	// send raw request data to html socket server
	sock.Write(dump)
	resp, err := ioutil.ReadAll(sock)
	defer sock.Close()

	conn, respBuf, err := c.Writer.Hijack()
	defer conn.Close()

	checkErr(err)
	respBuf.Write(resp)
	respBuf.Flush()

}

func checkErr(err error) {

	if err != nil {
		log.Println(err)
	}
}
