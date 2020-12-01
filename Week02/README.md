# error

Go error 就是一个普通的接口。该接口定义在builtin/中
``go
type error interface {
    Error() string
}
``
Error 特性

- 简单
- 考虑失败，而不是成功（Plan for  failure， not success）
- 没有隐藏的控制流
- 完全交给调用者控制error
- Error are values
## Error Type
### Sentinel Error
定义：预定义的错误
缺点： 不够灵活，缺少上下文，包之间产生以来
结论 ：尽量避免使用
### Error types 
定义：实现了了error接口的自定义类型。比如
``go
type MyError struct {
    Msg string
    File string
    Line int
}

func (e *MyError) Error() string{
    return fmt.Sprintf("%s:%d:%s", e.File, e.Line, e.Msg)
}
func test() error {
    return &MyError{"Something happened", "server.go", 52}
}
``
缺点：调用者需要使用类型断言和类型switch，就需要将自定义Error变为public。导致和调用者产生强耦合，从而导致api变得脆弱
优点：能够包装底层错误以及提供更多的上线问， 比如os.PathError
### Opaque errros
定义： 不透明错误处理
缺点：
优点：

- 对上层调用透明，只返回错误而不假设其内容
- 断言错误实现的具体行为
``go
type temporary interface {
    Temporary() bool
}
func IsTemporaray(err error) bool {
    te, ok := err.(temporary)
    return ok && te.Temporary()
}
``


## Wrap errors
you should only handle errors once. Handling an error means inspecting the error value, and making a single decision.


``go
func ReadFile(path string) ([]byte, error) {
    ...
    return nil, erros.Wrap(err, "open failed")
}

func main() {
    _, err := ReadFile()
    fmt.Printf("original err:%T %v\n", errors.Cause(err), errors.Cause(err))
}
``
如果希望在error中携带一些上下文信息，可以使用 `WithMessage` 
``go
...

return errors.WithMessage(err, "Could not read Config")
``
### 最佳实践：
引入github.com/pkg/errors

- 如果和其他库进行写作，考虑使用errors.Wrap或者errosWrapf保存堆栈信息。
- 直接返回错误，而不是将错误导出打印
- 在顶层程序或者工作的goroutine顶部（请求入口），使用`%+v`来打印详细堆栈信息
- 使用errors.Cause获取root error
## GO 1.13
go1.13为errors和fmt引入了新特性， 以简化处理包含其他错误的错误。
errors.UnWrap 以及两个用于检查错误的新函数 `Is` `As` ，使用fmt.Errorf想错误添加附加信息



# 作业

```go

type User struct {
    userId   int
    userName string
}

type Dao struct {
}

type Service struct {
}

var RecordNotFound = errors.New("record not found")

var dao = &Dao{}
var service = &Service{}

func main() {
    user, err := service.getUser(1)
    if errors.Is(err, RecordNotFound) {
            log.Printf("用户%d不存在,err:%+v\n", 1, err)
        }
}

func (d *Dao) getUser(userId int) (*User, error) {

    err := sql.ErrNoRows

    if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.Wrap(RecordNotFound, "dao.GetUser")
        }

    return nil, nil
}

func (s *Service) getUser(userId int) (*User, error) {
    user, err := dao.getUser(userId)

    if err != nil {
            return nil, errors.Wrap(err, "service.getUser")
        }
}

```
