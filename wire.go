package wire

import (
	"os"

	svc "github.com/judwhite/go-svc"
)

//wire 模块用来管理业务的服务依赖与服务启停
var wire *Wire

// init 初始化wire
func init() {
	wire = &Wire{
		sequence: []Service{},
	}
}

// Wire record service dependence
type Wire struct {
	sequence []Service
}

// Message struct for notify
type Message struct {
	Key  string
	Data interface{}
}

// Service  定义服务接口要求
type Service interface {
	// 初始化服务
	Init() error
	// 启动服务
	Start() error
	// 停止服务
	Stop() error
	// Notify 发送通知消息
	Notify(Message) error
}

// Init init service
func (w *Wire) Init(env svc.Environment) error {
	var err error
	for _, svc := range w.sequence {
		err = svc.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

//Start call when start
func (w *Wire) Start() error {
	var err error
	for _, svc := range w.sequence {
		err = svc.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

//Notify call when change
func (w *Wire) Notify(msg Message) error {
	var err error
	for _, svc := range w.sequence {
		err = svc.Notify(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop call when stop
func (w *Wire) Stop() error {
	var err error
	// reverse access
	for i := len(w.sequence); i > 0; i-- {
		svc := w.sequence[i-1]
		err = svc.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

type environment struct{}

func (environment) IsWindowsService() bool {
	return false
}

// Append append service to queue
func (w *Wire) Append(service Service) {
	w.sequence = append(w.sequence, service)
}

// Init   init  all service
func Init() error {
	env := environment{}
	return wire.Init(env)
}

// Start  start all service
func Start() error {
	return wire.Start()
}

// Stop stop all service
func Stop() error {
	return wire.Stop()
}

//Append append service to queue
func Append(service Service) {
	wire.Append(service)
}

//Notify broadcast  message to all service
func Notify(msg Message) error {
	return wire.Notify(msg)
}

// Run start run wire and monitor signal
func Run(sig ...os.Signal) error {
	var err error
	err = svc.Run(wire, sig...)
	return err
}

// BaseService base type for Service Interface
type BaseService struct {
}

// Init be called when call  wire.Init
func (bs BaseService) Init() error {
	return nil
}

// Start  be called when call  wire.Start
func (bs BaseService) Start() error {
	return nil
}

// Stop be called when call  wire.Stop
func (bs BaseService) Stop() error {
	return nil
}

// Notify be called when call  wire.Notify
func (bs BaseService) Notify(Message) error {
	return nil
}
