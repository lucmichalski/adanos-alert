package template


import (
	"strings"
	"testing"

	pkgJSON "github.com/mylxsw/adanos-alert/pkg/json"
	"github.com/mylxsw/go-toolkit/file"
)

var content = `
create mode 100644 .gitignore
create mode 100644 Gopkg.lock
create mode 100644 Gopkg.toml
create mode 100644 LICENSE
create mode 100644 README.md
create mode 100644 alert/dingding/dingding.go
create mode 100644 api/alert_rule.go
create mode 100644 api/api.go
create mode 100644 api/message.go
create mode 100644 message-view.tmpl
create mode 100644 parser/parser.go
create mode 100644 parser/parser_test.go
create mode 100644 route/api.go
create mode 100644 route/router.go
create mode 100644 route/web.go
create mode 100644 server/main.go
create mode 100644 storage/alert_rules.go
create mode 100644 storage/message.go
create mode 100644 storage/message_group.go
create mode 100644 storage/model.go
create mode 100644 storage/sequence.go
create mode 100644 storage/user.go
create mode 100644 web/message.go`

func TestCutOffFunc(t *testing.T) {
	if len(cutOff(40, content)) > 40 {
		t.Errorf("CutOff函数执行异常")
	}
}

func TestLeftIdent(t *testing.T) {
	for _, line := range strings.Split(leftIdent("....", content), "\n") {
		if !strings.HasPrefix(line, "....") {
			t.Errorf("leftIdent函数执行异常")
		}
	}
}

func TestJsonGet(t *testing.T) {
	var jsonContent = `{
    "message": "sms_send_failed",
    "context": {
        "msg": "短信发送失败，该错误不允许重试其它通道，请人工介入处理",
        "sms": {
            "id": 30139,
            "phone": "15923356841",
            "app_id": 1,
            "template_params": {
                "material_no": "JN027",
                "operator_name": "陈勇",
                "sence_name": "设备维护",
                "node_name": "故障报修",
                "material_name": "机械压机",
                "operate_time": "2018-11-01 23:12:41"
            },
            "template_id": 5,
            "status": 2,
            "created_at": "2018-11-01 23:12:41",
            "updated_at": "2018-11-03 23:31:01"
        },
        "final_channel": "亿美软通",
        "ack": {
            "msg": "停机或空号",
            "code": "MI:0013"
        },
        "line": "/data/webroot/message/app/Components/SmsHandler.php:383"
    },
    "level": 400,
    "level_name": "ERROR",
    "channel": "custom_cmd",
    "datetime": "2018-11-03 23:31:01",
    "extra": {
        "ref": "20181103233101-5bddbf35de697"
    }
}`

	if pkgJSON.Get("message", "", jsonContent) != "sms_send_failed" {
		t.Errorf("test failed")
	}

	if pkgJSON.Get("context.msg", "", jsonContent) != "短信发送失败，该错误不允许重试其它通道，请人工介入处理" {
		t.Errorf("test failed")
	}

}

func TestJsonGets(t *testing.T) {
	var jsonContent1 = `{
	"msg": "abc",
	"context": {
		"c1": { "user": 123 },
		"c2": { "user": 345 }
	}
}`
	if pkgJSON.Gets("msg,message", "", jsonContent1) != "abc" {
		t.Errorf("test failed")
	}

	var jsonContent2 = `{
    "message": "sms_send_failed"
}`
	if pkgJSON.Gets("msg,message", "", jsonContent2) != "sms_send_failed" {
		t.Errorf("test failed")
	}


}

func TestStartsWith(t *testing.T) {
	if !startsWith("stacktrace", "stack") {
		t.Errorf("test failed")
	}

	if !startsWith("stacktrace", "trace", "stack") {
		t.Errorf("test failed")
	}

	if startsWith("stackstrace", "trace") {
		t.Errorf("test failed")
	}
}

func TestEndsWith(t *testing.T) {
	if !endsWith("stacktrace", "trace") {
		t.Errorf("test failed")
	}

	if !endsWith("stacktrace", "stack", "trace") {
		t.Errorf("test failed")
	}

	if endsWith("stackstrace", "test") {
		t.Errorf("test failed")
	}
}

func TestParseMysqlSlowlog(t *testing.T) {

	testcases := map[string]string{
		"./test/mysql-slow.log":   "\nSQL: select count(*) as aggregate from `ms_sms_message_history` where exists (select * from `ms_sms_message` where `ms_sms_message_history`.`sms_id` = `ms_sms_message`.`id` and `phone` = '17623803369' and `created_at` >= '2019-06-26 00:00:00') and `channel_id` = 1\nDB: core_message\nElapse: 3.297862\nRowsExamined: 1915\nUser: core_message\n",
		"./test/mysql-slow-2.log": "\nSQL: INSERT INTO XXX(ID,client_id,client_name)\nSELECT * FROM XXX\nDB: bd_yunsombi\nElapse: 229.595303\nRowsExamined: 374611\nUser: bd_yunsombi\n",
		"./test/mysql-slow-3.log": "\nSQL: select child.*, parent.material_id as new_material_id,parent.material_ins_id as new_material_ins_id\n    from workflow_instance parent\n    join workflow_instance child on parent.id = child.parent_wf_inst_id\n    where child.material_ins_id  <>  parent.material_ins_id\nDB: \nElapse: 1.200693\nRowsExamined: 1041475\nUser: mat_lifecycle\n",
	}

	templateContent := `{{ $slog := mysql_slowlog .body }}
SQL: {{ $slog.sql }}
DB: {{ $slog.database }}
Elapse: {{ $slog.query_time }}
RowsExamined: {{ $slog.rows_examined }}
User: {{ $slog.user }}
`
	for f, tc := range testcases {
		mysqlSlowlog, err := file.FileGetContents(f)
		if err != nil {
			t.Errorf("test failed: %s", err)
		}

		res, err := Parse(templateContent, map[string]string{"body": mysqlSlowlog})
		if err != nil {
			t.Errorf("test failed: %s", err)
		}

		if res != tc {
			t.Errorf("test failed: %s", res)
		}
	}
}