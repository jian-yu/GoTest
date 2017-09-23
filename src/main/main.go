package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"studentinfomanagement/src/course"
	"studentinfomanagement/src/klass"
	"studentinfomanagement/src/lib"
	"studentinfomanagement/src/score"
	"studentinfomanagement/src/student"

	"github.com/astaxie/beego/orm"
)

var myOrm orm.Ormer

func init() {
	orm.RegisterDataBase("default", "mysql", "root:0926@/StudentInfoManagement?charset=utf8", 30)
	orm.RegisterModel(new(lib.Manager), new(lib.Student), new(lib.StudentPassword), new(lib.Class), new(lib.Score), new(lib.Course))
}

func main() {
	myOrm = orm.NewOrm()
	myOrm.Using("StudentInfoManagement")
	http.HandleFunc("/", manager)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println(err)
	}
}
func manager(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	fmt.Printf("Method is %s , Path is %s\n", method, path)
	if method == "POST" {
		switch path {
		case "/userlogin":
			userLogin(res, req)
		case "/managerlogin":
			managerLogin(res, req)
		case "/student/update":
			studentUpdate(res, req)
		case "/student/insert":
			studentInsert(res, req)
		case "/student/delete":
			studentDelete(res, req)
		case "/class/insert":
			classInsert(res, req)
		case "/class/delete":
			classDelete(res, req)
		case "/score/update":
			scoreUpdate(res, req)
		case "/score/delete":
			scoreDelete(res, req)
		case "/score/insert":
			scoreInsert(res, req)
		case "/course/delete":
			courseDelete(res, req)
		case "/course/insert":
			courseInsert(res, req)
		default:
			responseOp(401, "URL不存在", nil, res)
		}
	} else if method == "GET" {
		switch path {
		case "/student/query":
			studentDetail(res, req)
		case "/class/queryall":
			classQueryAll(res, req)
		case "/student/queryby":
			studentQueryByClassID(res, req)
		case "/student/queryall":
			studentQueryAll(res, req)
		case "/course/queryall":
			courseQueryAll(res, req)
		case "/score/queryall":
			scoreQueryAll(res, req)
		case "/score/queryby":
			scoreQueryBy(res, req)
		case "/score/queryclass":
			scoreQueryClass(res, req)
		default:
			responseOp(401, "URL不存在", nil, res)
		}
	} else {
		responseOp(401, "方法不存在", nil, res)
	}
}

//普通用户登录
func userLogin(res http.ResponseWriter, req *http.Request) {
	loginResult, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(loginResult))
	defer req.Body.Close()
	if err != nil {
		responseOp(401, err.Error(), nil, res)
	} else {
		var resultJSON lib.LoginRequest
		err := json.Unmarshal(loginResult, &resultJSON)
		fmt.Println(resultJSON)
		if err != nil {
			responseOp(401, err.Error(), nil, res)
		} else {
			password := resultJSON.Body.Password
			var stuPw lib.StudentPassword
			err := myOrm.QueryTable("student_password").Filter("id", resultJSON.Body.Id).One(&stuPw)
			if err == nil {
				if password == stuPw.Password {
					responseOp(210, "login successful!", nil, res)
				} else {
					responseOp(410, "login error! password error", nil, res)
				}
			} else {
				fmt.Println(err)
				responseOp(410, err.Error(), nil, res)
			}
		}
	}
}

//管理员登录
func managerLogin(res http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	fmt.Println(string(b))
	if err != nil {
		responseOp(402, "数据读取错误", nil, res)
	} else {
		var managerRes lib.ManagerRequest
		err := json.Unmarshal(b, &managerRes)
		if err != nil {
			responseOp(402, "数据解析错误", nil, res)
		} else {
			password := managerRes.Body.Password
			var manager lib.Manager
			err := myOrm.QueryTable("Manager").Filter("name", managerRes.Body.Name).One(&manager)
			if err != nil {
				responseOp(404, "该用户未注册", nil, res)
			} else {
				if password == manager.Password {
					responseOp(201, "Manager login successful!", manager, res)
				} else {
					responseOp(403, "Manager login password error", nil, res)
				}
			}
		}
	}

}
func classDelete(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	fmt.Println(string(jsonByte))
	if err != nil {
		responseOp(401, "数据读取错误", nil, res)
		return
	}
	var c lib.ClassRequest
	err = json.Unmarshal(jsonByte, &c)
	if err != nil {
		responseOp(401, "数据解析失败", nil, res)
		return
	}
	co := klass.NewClassOperator(myOrm)
	_, err = co.Delete(c.Body.Id)
	if err == nil {
		responseOp(210, "删除成功", nil, res)
	} else {
		responseOp(410, err.Error(), nil, res)
	}

}
func classInsert(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	fmt.Println(string(jsonByte))
	if err != nil {
		responseOp(401, "数据读取错误", nil, res)
		return
	}
	var c lib.ClassRequest
	err = json.Unmarshal(jsonByte, &c)
	if err != nil {
		responseOp(401, "数据解析失败", nil, res)
		return
	}
	co := klass.NewClassOperator(myOrm)
	klassAfterInsert, err := co.Insert(&c.Body)
	if err == nil {
		responseOp(210, "插入成功", klassAfterInsert, res)
	} else {
		responseOp(410, err.Error(), nil, res)
	}

}

func classQueryAll(res http.ResponseWriter, req *http.Request) {
	co := klass.NewClassOperator(myOrm)
	klass, err := co.QueryAll()
	fmt.Println(klass)
	if err != nil {
		responseOp(410, err.Error(), nil, res)
		return
	}
	responseOp(210, "查询成功", klass, res)
}

func studentDetail(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err == nil {
		idStr := req.Form.Get("id")
		sto, _ := student.NewStudentOperator(myOrm)
		id, _ := strconv.Atoi(idStr)
		stuAfterQuery, err := sto.Query(id)
		if err != nil {
			responseOp(410, err.Error(), nil, res)
			return
		}
		responseOp(210, "查询成功", stuAfterQuery, res)
	}

}
func studentUpdate(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(401, "数据读取错误", nil, res)
		return
	}
	var stu lib.StudentRequest
	err = json.Unmarshal(jsonByte, &stu)
	if err != nil {
		responseOp(401, "数据解析失败", nil, res)
		return
	}
	sto, _ := student.NewStudentOperator(myOrm)
	stuAfterUpdate, err := sto.Update(stu.Body.Id, stu.Body)
	if err != nil {
		responseOp(410, err.Error(), nil, res)
		return
	}
	responseOp(210, "修改成功", stuAfterUpdate, res)
}
func studentDelete(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(401, "数据读取错误", nil, res)
		return
	}
	var stu lib.StudentRequest
	err = json.Unmarshal(jsonByte, &stu)
	if err != nil {
		responseOp(401, "数据解析失败", nil, res)
		return
	}
	sto, _ := student.NewStudentOperator(myOrm)
	_, err = sto.Delete(stu.Body.Id)
	if err == nil {
		responseOp(210, "删除成功", nil, res)
	} else {
		responseOp(410, err.Error(), nil, res)
	}
}
func studentInsert(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(402, err.Error(), nil, res)
		return
	}
	var stu lib.StudentRequest
	err = json.Unmarshal(jsonByte, &stu)
	fmt.Println(stu)
	if err != nil {
		responseOp(402, "数据解析错误", nil, res)
		return
	}
	sto, _ := student.NewStudentOperator(myOrm)
	ss, err := sto.Insert(&stu.Body)
	if err != nil {
		responseOp(410, err.Error(), nil, res)
		return
	}
	responseOp(210, "插入成功", ss, res)
}

func studentQueryByClassID(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err == nil {
		idStr := req.Form.Get("class_id")
		classID, _ := strconv.Atoi(idStr)
		sto, _ := student.NewStudentOperator(myOrm)
		stuArr, err := sto.QueryStudentByClassID(classID)
		if err == nil {
			responseOp(210, "查询成功", stuArr, res)
			return
		}
		responseOp(410, err.Error(), nil, res)
	}
}

func studentQueryAll(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err == nil {
		sto, _ := student.NewStudentOperator(myOrm)
		stuArr, err := sto.QueryStudentAll()
		if err == nil {
			responseOp(210, "查询成功", stuArr, res)
			return
		}
		responseOp(410, err.Error(), nil, res)
	}
}

func scoreUpdate(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(401, "数据读取失败", nil, res)
		return
	}
	var sc lib.ScoreRequest
	err = json.Unmarshal(jsonByte, &sc)
	fmt.Println(sc)
	if err != nil {
		responseOp(401, "数据解析失败", nil, res)
		return
	}
	sco := score.NewScoreOparetor(myOrm)
	scAfterUpdate, err := sco.Update(sc.Body.Student_id, sc.Body)
	if err != nil {
		responseOp(410, err.Error(), nil, res)
		return
	}
	responseOp(210, "更新成功", scAfterUpdate, res)
}

func scoreInsert(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(401, "数据读取失败", nil, res)
		return
	}
	var sc lib.ScoreRequest
	err = json.Unmarshal(jsonByte, &sc)
	fmt.Println(sc)
	if err != nil {
		responseOp(401, "数据解析失败", nil, res)
		return
	}
	sco := score.NewScoreOparetor(myOrm)
	scoAfterInsert, err := sco.Insert(&sc.Body)
	if err == nil {
		responseOp(210, "插入成功", scoAfterInsert, res)
	} else {
		responseOp(401, err.Error(), nil, res)
	}

}

func scoreDelete(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(401, "数据读取失败", nil, res)
		return
	}
	var sc lib.ScoreRequest
	err = json.Unmarshal(jsonByte, &sc)
	fmt.Println(sc)
	if err != nil {
		responseOp(401, "数据解析失败", nil, res)
		return
	}
	sco := score.NewScoreOparetor(myOrm)
	_, err = sco.Delete(sc.Body.Student_id, sc.Body.Course_id)
	if err != nil {
		responseOp(410, err.Error(), nil, res)
		return
	}
	responseOp(210, "删除成功", nil, res)

}

func scoreQueryBy(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err == nil {
		SIDStr := req.Form.Get("student_id")
		SID, _ := strconv.Atoi(SIDStr)
		sco := score.NewScoreOparetor(myOrm)
		scoreAfterQuery, err := sco.QueryBy(SID)
		if err == nil {
			responseOp(210, "查询成功", scoreAfterQuery, res)
			return
		}
		responseOp(410, err.Error(), nil, res)
	} else {
		responseOp(401, "数据解析错误", nil, res)

	}
}

func scoreQueryAll(res http.ResponseWriter, req *http.Request) {
	sco := score.NewScoreOparetor(myOrm)
	scoreAfterQuery, err := sco.QueryAll()
	if err == nil {
		responseOp(210, "查询成功", scoreAfterQuery, res)
		return
	}
	responseOp(410, err.Error(), nil, res)

}

func scoreQueryClass(res http.ResponseWriter, req *http.Request) {
	fmt.Println("scoreQueryClass")
	err := req.ParseForm()
	if err == nil {
		CIDStr := req.Form.Get("class_id")
		CID, _ := strconv.Atoi(CIDStr)
		fmt.Println(CID)
		so := score.NewScoreOparetor(myOrm)
		allScore, err := so.QueryClass(CID)
		if err == nil {
			responseOp(210, "查询成功", allScore, res)
		} else {
			responseOp(410, err.Error(), nil, res)
		}
	} else {
		responseOp(401, "数据解析错误", nil, res)
	}
}

func courseInsert(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(402, err.Error(), nil, res)
		return
	}
	var cr lib.CourseRequset
	err = json.Unmarshal(jsonByte, &cr)
	fmt.Println(cr)
	if err != nil {
		responseOp(402, "数据解析错误", nil, res)
		return
	}
	co := course.NewCourseOperator(myOrm)
	coAfterInsert, err := co.Insert(&cr.Body)
	if err == nil {
		responseOp(210, "插入成功", coAfterInsert, res)
		return
	}
	responseOp(410, err.Error(), nil, res)
}

func courseDelete(res http.ResponseWriter, req *http.Request) {
	jsonByte, err := dataParse(req)
	if err != nil {
		responseOp(402, err.Error(), nil, res)
		return
	}
	var cr lib.CourseRequset
	err = json.Unmarshal(jsonByte, &cr)
	fmt.Println(cr)
	if err != nil {
		responseOp(402, "数据解析错误", nil, res)
		return
	}
	co := course.NewCourseOperator(myOrm)
	coAfterInsert, err := co.Delete(&cr.Body)
	if err == nil {
		responseOp(210, "删除成功", coAfterInsert, res)
		return
	}
	responseOp(410, err.Error(), nil, res)
}
func courseQueryAll(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err == nil {
		co := course.NewCourseOperator(myOrm)
		cosAfterInsert, err := co.QueryAll()
		if err == nil {
			responseOp(210, "查询成功", cosAfterInsert, res)
			return
		}
		responseOp(410, err.Error(), nil, res)
	}
	responseOp(410, "请输入正确的URL", nil, res)
}

func dataParse(req *http.Request) (jsonByte []byte, err error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	fmt.Println(string(b))
	if err != nil {
		return nil, errors.New("数据解析错误")
	}
	return b, nil
}

func responseOp(status int, message string, body interface{}, res http.ResponseWriter) {
	response := lib.Response{Status: status, Message: message, Body: body}
	resBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	res.Write(resBytes)
}
