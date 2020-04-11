
<p align="center">
<img src="https://travis-ci.com/dashjay/ustb_sso_server.svg?token=bpayQtpryFqAhnyVqjxy&branch=master" alt="">
</p>


# 关于 SSO-Server

这个项目可以接通北京科技大学微校园服务，提供简易的用户认证API，使用非常简单，针对了所有贝壳校园的API，完成以下服务：

- [x] grade接口【获取用户成绩】
- [x] course接口【获取用户课表】
- [x] cet四六级接口 【获取用户46级成绩】
- [x] head头像【获取用户在贝壳校园的头像】
- [x] info 【获取用户信息】
- [x] exam 【获取用户考试安排信息】
- [ ] infomore 【更多信息】用户自主录入的信息

# 接口说明
> http 接口

### 认证接口`/auth`
请求 `/auth` 时必须携带`fuck`头，模拟`curl`请求如下

```bash
curl "localhost:80/auth?union_id=1" -H 'fuck:fuck'
# 返回 {"code":101,"url":"","uid":"xxxxxxxx","msg":""} json结构在下方有解释
```

```gotemplate
// 返回结构定义
type AuthStruct struct {
	// Code 状态参数：如果 != 1 说明有数据传出，
	// Code 1说明有错读msg即可
	// Code 101已认证 URL 为空 UID 为学号
	// Code 102未认证 URL 为认证URL UID 为空
	Code int `json:"code" bson:"code"`

	//  "https://sis.ustb.edu.cn/auth" 开头的，传给学生，点开之后可以完成认证
	URL string `json:"url" bson:"url"`

	// 用户学号
	UID string `json:"uid" bson:"uid"`

	// 信息
	Msg string `json:"msg" bson:"msg"`
}

```

### 函数接口/func
请求`/func`时也必须携带fuck头，模拟curl的几个请求如下

#### 请求成绩
`curl "localhost:80/func?union_id=1&func=grade" -H 'fuck: fuck'`

#### 返回成绩

```gotemplate
{
  "msg" : "OK",
  "body" : [ {
    "schoolYear" : 2017,
    "semester" : 1,
    "courseNo" : "10601831",
    "courseName" : "工科数学分析I",
    "period" : 0,
    "credit" : 6.0,
    "examResults" : 72.0,
    "originResults" : 0.0
  }{
    ......
    }],
  "code" : "SUCCESS"
}
```

#### 请求课表

`curl "localhost:80/func?union_id=1&func=course" -H 'fuck: fuck'`

```gotemplate
{
  "msg" : "OK",
  "body" : {
    "map" : {
      "1" : { //第几节课
        "1" : [ ], // 星期1
        "2" : [ ], // 星期2
        "3" : [ {
          "podium.podiumName" : "人工智能",
          "SKZCZFC" : "1-6周",
          "weeks" : "222222000000000000000000000000",
          "section" : 1,
          "userName" : "41724235",
          "classroom.id" : "3334191",
          "courseName" : "人工智能",
          "classroom.roomNickname" : "逸304",
          "dayOfWeek" : 3,
          "podium.courseName" : "人工智能",
          "currentWeek" : "Y",
          "id" : 7,
          "podium.id" : 1002
        }
     .....
        "7" : [ ]
      }
    },
    "NO_SECTION" : [ {
      "classroom.id" : null,
      "podium.podiumName" : "软件工程课程设计(实验)",
      "SKZCZFC" : null,
      "courseName" : "软件工程课程设计(实验)",
      "classroom.roomNickname" : null,
      "weeks" : null,
      "dayOfWeek" : null,
      "podium.courseName" : "软件工程课程设计(实验)",
      "section" : null,
      "id" : 16,
      "userName" : "41724235",
      "podium.id" : null
    }, {
        ...
    } ]
  },
  "code" : "SUCCESS"
}
```

#### 其它请求我不再逐个示范

| funcname | explain      |
| -------- | ------------ |
| grade    | 获取成绩     |
| course   | 获取课表     |
| cet      | 获取46级     |
| head     | 获取头像     |
| info     | 获取用户信息 |
| ....     | ....         |

# dockerize

在docker化方面我使用了先在`golang:alpine`中编译，然后拷贝到apline中运行，生成的`docker image`只有20M大小，docker真的是好东西。

```
FROM golang:alpine AS build
.....
COPY ./ /go/sso_server
RUN go build -o sso_server ./main.go
FROM alpine:latest
...
COPY --from=build /go/sso_server /opt/
EXPOSE 80/tcp 81/tcp
ENTRYPOINT ["./sso_server"]
```



直接在目录下跑

```
docker build . -t sso_server 
docker run -d -p 1080:80 -p 1081:81 sso_server:latest
```

太棒了，现在使用它只需要服务器20M的空间啦~