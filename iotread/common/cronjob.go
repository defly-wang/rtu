package common

// /定义一个类型 包含一个int类型参数和函数体
type funcIntJob struct {
	num      int
	function func(int) bool
}

// 实现这个类型的Run()方法 使得可以传入Job接口
// 在次函数中处理
func (my *funcIntJob) Run() {

	if nil != my.function {
		//ret :=
		my.function(my.num)
		/*
			{
				//fmt.Println(ret)
				return
			}
		*/
	}

}

// 非必须  返回一个urlServeJob指针
func FuncIntJob(num int, function func(int) bool) *funcIntJob {

	instance := &funcIntJob{
		num:      num,
		function: function,
	}
	return instance
}
