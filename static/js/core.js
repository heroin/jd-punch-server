Date.prototype.format = function(format) {

    var o = {
        "M+" :this.getMonth() + 1,
        "d+" :this.getDate(),
        "h+" :this.getHours(),
        "m+" :this.getMinutes(),
        "s+" :this.getSeconds(),
        "q+" :Math.floor((this.getMonth() + 3) / 3),
        "S" :this.getMilliseconds()
    }

    if (/(y+)/.test(format)) {
        format = format.replace(RegExp.$1, (this.getFullYear() + "")
                .substr(4 - RegExp.$1.length));
    }

    for (var k in o) {
        if (new RegExp("(" + k + ")").test(format)) {
            format = format.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k]
                    : ("00" + o[k]).substr(("" + o[k]).length));
        }
    }
    return format;
}

function make_task_item(username, password, trigger, start) {
    var time = new Date(trigger * 1000);
    html =  "                    <tr username=\"" + username + "\" trigger=\"" + trigger + "\">\n";
    html += "                      <td>" + username + "</td>\n";
    html += "                      <td>" + password + "</td>\n";
    html += "                      <td class=\"center\">" + time.format("yyyy-MM-dd hh:mm:ss") + "</td>\n";
    html += "                      <td class=\"center\">" + start + "</td>\n";
    html += "                      <td class=\"center\">\n";
    html += "                        <a class=\"task-cancel\" title=\"cancel\" href=\"javascript:void(0);\"><i class=\"icon-edit\"></i></a>\n";
    html += "                        <a class=\"task-remove\" title=\"remove\" href=\"javascript:void(0);\"><i class=\"icon-remove\"></i></a>\n";
    html += "                      </td>\n";
    html += "                    </tr>\n";
    return html;
}

function make_cancel_item(username, trigger) {
    var time = new Date(trigger * 1000);
    key = username + "-" + trigger;
    html =  "                    <tr username=\"" + username + "\" trigger=\"" + trigger + "\" key=\"" + key + "\">\n";
    html += "                      <td>" + key + "</td>\n";
    html += "                      <td class=\"center\">" + time.format("yyyy-MM-dd hh:mm:ss") + "</td>\n";
    html += "                      <td class=\"center\">\n";
    html += "                        <a class=\"cancel-remove\" title=\"remove\" href=\"javascript:void(0);\"><i class=\"icon-remove\"></i></a>\n";
    html += "                      </td>\n";
    html += "                    </tr>\n";
    return html;
}

function bind() {
    $("#tab-task-list .task-cancel").click(function () {
        username = $.trim($(this).parent().parent().attr("username"));
        trigger = $.trim($(this).parent().parent().attr("trigger"));
        if (username != "" && trigger != "") {
            $.post("/cancel/add", {"username": username, "trigger": trigger}, function (result) {
                if (result == "success") {
                    load();
                }
            });
        }
    });

    $("#tab-task-list .task-remove").click(function () {
        username = $.trim($(this).parent().parent().attr("username"));
        trigger = $.trim($(this).parent().parent().attr("trigger"));
        if (username != "" && trigger != "") {
            $.post("/user/del", {"username": username, "trigger": trigger}, function (result) {
                if (result == "success") {
                    load();
                }
            });
        }
    });

    $("#tab-cancel-list .cancel-remove").click(function () {
        username = $.trim($(this).parent().parent().attr("username"));
        trigger = $.trim($(this).parent().parent().attr("trigger"));
        if (username != "" && trigger != "") {
            $.post("/cancel/del", {"username": username, "trigger": trigger}, function (result) {
                if (result == "success") {
                    load();
                }
            });
        }
    });
}

function load() {
    $.get("/task.json", function(data) {
        $("#tab-content #tab-task-list tbody tr").remove();
        result = eval("(" + data + ")");
        users = result["users"];
        for (i in users) {
            $("#tab-content #tab-task-list tbody").append(make_task_item(users[i].username, users[i].password, users[i].trigger, users[i].start));
        }
        cancel = result["cancel"];
        $("#tab-content #tab-cancel-list tbody tr").remove();
        for (i in cancel) {
            data_cancel = cancel[i].split("-");
            $("#tab-content #tab-cancel-list tbody").append(make_cancel_item(data_cancel[0], data_cancel[1]));
        }
        bind();
    });
}

$(document).ready(function () {
    $("#trigger-widget").datetimepicker({
        language: "us"
    });

    load();

    $("#btn-add-task").click(function () {
        username_var = $.trim($("#username").val());
        password_var = $.trim($("#password").val());
        start_var = $('input:radio[name="start"]:checked').val();
        trigger_var = $.trim($("#trigger").val());
        if (username_var != "" && password_var != "" && trigger_var != "") {
            trigger_date = new Date(trigger_var.substring(0, 4), trigger_var.substring(5, 7) - 1, trigger_var.substring(8, 10), trigger_var.substring(11, 13), trigger_var.substring(14, 16), trigger_var.substring(17, 19));
            trigger = trigger_date.getTime() / 1000;
            $.post("/user/add", {"username": username_var, "password": password_var, "trigger": trigger, "start": start_var}, function (result) {
                if (result == "success") {
                    load();
                    $("#add-task").removeClass("in active");
                    $("#nav-add-task").removeClass("active");
                    $("#task-list").addClass("in active");
                    $("#nav-task-list").addClass("active");
                }
            });
        }
    });
});
