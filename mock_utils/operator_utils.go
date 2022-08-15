package mock_utils

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/tjlcast/go_common/executor_utils"
	"github.com/tjlcast/go_common/file_utils"
	"github.com/tjlcast/go_common/log_utils"
	"github.com/tjlcast/go_common/rand_utils"
	"github.com/tjlcast/go_common/time_utils"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	OP_SUFFIX          = "op"
	OP_POINTER         = "."
	OP_ARGS_SPLITER    = "-"
	OP_COMMAND_SPLITER = "&"
)

type Command struct {
	Name        string
	CommandType int // 0: once 1: loop
	handler     func(string, string)
}

func GeneratCommandMapper(funcArr []func(string, string)) map[string]*Command {
	commandMapper := make(map[string]*Command)

	commands := []*Command{}
	for _, funcE := range funcArr {
		name := GetFuncName(funcE)
		commands = append(commands, &Command{
			Name:        name,
			CommandType: 0,
			handler:     funcE,
		})
	}
	for _, command := range commands {
		commandMapper[command.Name] = command
	}
	return commandMapper
}

func GetFuncName(function interface{}) string {
	originName := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	name := strings.Split(originName, ".")[1]
	return name
}

type MockOPClient struct {
	Name string

	commands        map[string]*Command
	clearCommanName string

	pool       *executor_utils.Pool
	taskMapper map[string][]*executor_utils.Task

	operatorParseChan chan string
	taskClearChan     chan string
}

func NewMockOPClient(name string, clearCommand string, commands map[string]*Command) *MockOPClient {
	client := &MockOPClient{Name: name}
	client.commands = commands

	if _, ok := client.commands[clearCommand]; !ok {
		panic("clearCommand is not in commands.")
	}
	client.clearCommanName = clearCommand
	return client
}

func (c *MockOPClient) PrintSample() {
	for commandName, _ := range c.commands {
		sprintf := fmt.Sprintf("touch %s-id-arg1.op", commandName)
		fmt.Println(sprintf)
	}
}

func (c *MockOPClient) clear() {
	for entity := range c.taskClearChan {
		entityTasks, ok := c.taskMapper[entity]
		if ok {
			delete(c.taskMapper, entity)
			for _, task := range entityTasks {
				task.Interupt = true
			}
		}
	}
}

func (c *MockOPClient) Loop(basePath string) {
	log_utils.Logger.Info(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	pool, _ := executor_utils.NewPool(32)
	c.pool = pool
	taskMap := make(map[string][]*executor_utils.Task)
	c.taskMapper = taskMap
	OperChan := make(chan string)
	c.operatorParseChan = OperChan
	c.taskClearChan = make(chan string)

	go c.clear()

	go func() {
		for rePath := range OperChan {

			rePath = strings.ReplaceAll(rePath, " ", "")
			rePath = strings.ReplaceAll(rePath, "./", "")

			op := strings.Split(file_utils.GetFileNameInPath(rePath), OP_POINTER)[0]

			flags := strings.Split(op, OP_ARGS_SPLITER)
			commandName := flags[0]

			var entityId string
			if len(flags) > 1 {
				entityId = flags[1]
			}

			var entityArgs string
			if len(flags) > 2 {
				entityArgs = flags[2]
			}

			log_utils.Logger.Warning("Get op: ", rePath)
			log_utils.Logger.Warning("Get command: ", commandName)

			command, ok := c.commands[commandName]
			if !ok {
				log_utils.Logger.Error("Not support: ", commandName)
				continue
			}

			if command != nil {
				log_utils.Logger.Warnf("Command: %s Entity: %s Args: %s\n", commandName, entityId, entityArgs)

				task := &executor_utils.Task{}
				task.Id = xid.New().String()
				task.Handler = func(v ...interface{}) {
					defer func() {
						task.Interupt = true
					}()
					log_utils.Logger.Info(">>>Submit a task -> " + entityId)
					if command.CommandType == 0 {
						command.handler(entityId, entityArgs)
					} else if command.CommandType == 1 {
						for !task.Interupt {
							command.handler(entityId, entityArgs)
							RandSleep(1)
						}
					} else {
						panic("Unknow commandType: " + strconv.Itoa(command.CommandType))
					}
					log_utils.Logger.Info(">>>Exit a task -> " + entityId)
				}

				err := c.pool.Put(task)
				if err != nil {
					log_utils.Logger.Error(err.Error())
				}
				tasks := taskMap[entityId]
				tasks = append([]*executor_utils.Task{}, tasks...)
				tasks = append(tasks, task)
				taskMap[entityId] = tasks

				if commandName == c.clearCommanName {
					c.taskClearChan <- entityId
				}
			}
		}
	}()

	go func() {
		for true {
			time.Sleep(1 * time.Second)
			commandFiles := file_utils.GetSuffixPaths(".", "op")
			if len(commandFiles) == 0 {
				continue
			}
			commands := strings.Split(commandFiles[0], "&")
			for _, command := range commands {
				if !strings.HasSuffix(command, "op") {
					continue
				}
				OperChan <- command
			}
			_ = os.Remove(commandFiles[0])
		}
	}()

	log_utils.Logger.Info("Test start>>>>>>>>>")
	loop()
}

func loop() {
	select {}
}

func RandSleep(sec int) {
	if rand_utils.GenRandInt(2) == 0 {
		time_utils.WaitSeconds(sec)
	}
}
