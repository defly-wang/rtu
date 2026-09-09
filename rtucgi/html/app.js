/* RTU 设备管理 - 前端逻辑 */
(function () {
  "use strict";

  var $ = function (id) { return document.getElementById(id); };
  function val(id) { return $(id).value.trim(); }
  function setVal(id, v) { $(id).value = (v === null || v === undefined) ? "" : v; }
  function showMsg(id, text, type) {
    var el = $(id);
    el.className = "msg " + (type || "info");
    el.textContent = text;
    el.classList.remove("hidden");
  }
function hideMsg(id) {
    var el = $(id);
    el.classList.add("hidden");
    el.textContent = "";
  }
  function showToast(text, type) {
    var el = $("toast");
    el.textContent = text;
    el.className = "toast " + (type || "info");
    el.classList.remove("hidden");
    var close = document.createElement("span");
    close.className = "toast-close";
    close.textContent = "×";
    close.addEventListener("click", function () { el.classList.add("hidden"); });
    el.appendChild(close);
    clearTimeout(showToast._t);
    showToast._t = setTimeout(function () { el.classList.add("hidden"); }, 5000);
  }
  function showGlobal(text, type) { showToast(text, type); }

  /* ================= 登录 ================= */
  function doLogin() {
    var u = val("login-user"), p = val("login-pass");
    hideMsg("login-msg");
    if (!u || !p) { showMsg("login-msg", "请输入用户名和密码", "error"); return; }
    var btn = $("btn-login");
    btn.disabled = true; btn.textContent = "登录中...";
    API.login(u, p).then(function () {
      btn.disabled = false; btn.textContent = "登 录";
      enterMain();
    }).catch(function (e) {
      btn.disabled = false; btn.textContent = "登 录";
      showMsg("login-msg", e.message, "error");
    });
  }

  function doLogout() {
    API.clearToken();
    location.hash = "";
    $("main-page").classList.add("hidden");
    $("login-page").classList.remove("hidden");
  }

  function enterMain() {
    $("login-page").classList.add("hidden");
    $("main-page").classList.remove("hidden");
    navigate("info");
  }

  /* ================= 导航 ================= */
  var currentPage = "info";
  function navigate(page) {
    currentPage = page;
    document.querySelectorAll("#nav .nav-item").forEach(function (a) {
      a.classList.toggle("active", a.dataset.page === page);
    });
    document.querySelectorAll(".content .page").forEach(function (p) {
      p.classList.toggle("hidden", p.id !== "page-" + page);
    });
    loadPage(page);
  }

  function loadPage(page) {
    switch (page) {
      case "info": loadInfo(); break;
      case "status": loadStatus(); break;
      case "base": loadBase(); break;
      case "mqtt": loadMqtt(); break;
      case "webapi": loadWebapi(); break;
      case "iot": loadIotList(); break;
      case "power": break; /* 由按钮触发 */
      case "his": loadHis(); break;
      case "test": break;
    }
  }

  /* ================= 系统信息 ================= */
  var sysinfoLabels = {
    read_version: "采集程序版本", rtu_version: "客户端程序版本", cgi_version: "Web 管理版本",
    version: "系统版本", uptime: "运行时间 (秒)", memtotal: "内存总量", memfree: "空闲内存",
    disktotal: "磁盘总量", diskfree: "磁盘剩余", mac: "MAC 地址", ip: "IP 地址",
    gprs_ip: "4G 内网 IP", gprs_imsi: "4G IMSI", gprs_ccid: "4G CCID", gprs_csq: "4G 信号强度"
  };
  function loadInfo() {
    API.getInfo().then(function (d) {
      var html = "";
      for (var k in sysinfoLabels) {
        if (d.data && d.data[k] !== undefined) html += kvItem(sysinfoLabels[k], d.data[k]);
      }
      $("info-body").innerHTML = html;
    }).catch(function (e) { $("info-body").innerHTML = '<div class="msg error">' + esc(e.message) + "</div>"; });
  }

  /* ================= 运行状态 ================= */
  var statusLabels = {
    starttime: "启动时间", count: "采集次数", readcount: "读取成功次数",
    mqttcount: "MQTT 发送次数", webapicount: "WebAPI 发送次数",
    runing: "运行状态", version: "程序版本"
  };
  function fmtStatus(d) {
    var copy = {};
    for (var k in d) {
      if (k === "runing") copy[k] = d[k] ? "运行中" : "未运行";
      else copy[k] = d[k];
    }
    return copy;
  }
  function loadStatus() {
    API.getStatus().then(function (d) {
      var html = "";
      var s = d.data || {};
      var f = fmtStatus(s);
      for (var k in statusLabels) {
        if (f[k] !== undefined) html += kvItem(statusLabels[k], f[k]);
      }
      $("status-body").innerHTML = html;
    }).catch(function (e) { $("status-body").innerHTML = '<div class="msg error">' + esc(e.message) + "</div>"; });
  }

  /* ================= 基础设置 ================= */
  function loadBase() {
    API.getBase().then(function (d) {
      var b = d.data || {};
      setVal("base-id", b.id);
      setVal("base-uint", b.uint);
      setVal("base-iotnumber", b.iotnumber);
      setVal("base-datadir", b.datadir);
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }
  function saveBase() {
    var d = {
      id: val("base-id"), uint: val("base-uint"),
      iotnumber: parseInt(val("base-iotnumber"), 10) || 0,
      datadir: val("base-datadir")
    };
    return API.saveBase(d);
  }

  /* ================= MQTT 设置 ================= */
  function loadMqtt() {
    API.getMqtt().then(function (d) {
      var m = d.data || {};
      setVal("mqtt-used", m.used ? "true" : "false");
      setVal("mqtt-host", m.host);
      setVal("mqtt-port", m.port || "");
      setVal("mqtt-topic", m.pubtoptic);
      setVal("mqtt-username", m.username);
      setVal("mqtt-password", m.password);
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }
  function saveMqtt() {
    var d = {
      host: val("mqtt-host"), port: parseInt(val("mqtt-port"), 10) || 0,
      pubtoptic: val("mqtt-topic"), username: val("mqtt-username"),
      password: val("mqtt-password"), used: $("mqtt-used").value === "true"
    };
    return API.saveMqtt(d);
  }

  /* ================= WebAPI 设置 ================= */
  function loadWebapi() {
    API.getWebapi().then(function (d) {
      var w = d.data || {};
      setVal("webapi-used", w.used ? "true" : "false");
      setVal("webapi-url", w.url);
      setVal("webapi-token", w.token);
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }
  function saveWebapi() {
    var d = {
      url: val("webapi-url"), token: val("webapi-token"),
      used: $("webapi-used").value === "true"
    };
    return API.saveWebapi(d);
  }

  /* ================= 传感器管理 ================= */
  var iotEditMode = "add"; /* add | edit */
  var iotEditIndex = 0;

  function loadIotList() {
    API.getIots().then(function (d) {
      var list = d.data || [];
      var rows = "";
      list.forEach(function (iot) {
        rows += "<tr>" +
          "<td>" + iot.iot + "</td>" +
          "<td>" + esc(iot.id) + "</td>" +
          "<td>" + esc(iot.type) + "</td>" +
          "<td>ttyS" + esc(iot.com) + "</td>" +
          "<td class='mono'>" + esc(iot.cron) + "</td>" +
          "<td>" + (iot.crc ? "已配置" : "-") + "</td>" +
          "<td>" +
            "<button class='btn btn-secondary btn-sm' onclick='App.editIotFromList(" + iot.iot + ")'>编辑</button> " +
            "<button class='btn btn-danger btn-sm' onclick='App.delIotFromList(" + iot.iot + ")'>删除</button>" +
          "</td></tr>";
      });
      $("iot-list").querySelector("tbody").innerHTML = rows;
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }

  function openIotForm(mode, iot) {
    iotEditMode = mode;
    $("iot-edit-title").textContent = mode === "edit" ? "编辑传感器 #" + iot.iot : "新增传感器";
    setVal("iot-iot", iot.iot == null ? "" : iot.iot);
    setVal("iot-id", iot.id || "");
    setVal("iot-type", iot.type || "");
    setVal("iot-cron", iot.cron || "0 * * * * *");
    setVal("iot-com", iot.com == null ? "" : iot.com);
    setVal("iot-baudrate", iot.baudrate || 9600);
    setVal("iot-databits", iot.databits || 8);
    setVal("iot-stopbits", iot.stopbits || 1);
    setVal("iot-paritymode", iot.paritymode == null ? "0" : iot.paritymode);
    setVal("iot-readbuff", iot.read_buff || "");
    setVal("iot-readlen", iot.read_len || 0);
    setVal("iot-readdelay", iot.read_delay || 0);
    setVal("iot-reciven", iot.recive_len || 0);
    setVal("iot-revivedatalen", iot.revive_data_len || 0);
    setVal("iot-revivedata", (iot.revive_data || []).join(","));
    setVal("iot-offset", iot.offset ? "true" : "false");
    setVal("iot-ratio", iot.ratio == null ? "" : iot.ratio);
    setVal("iot-crc", iot.crc || "CRC16");
    hideMsg("iot-edit-msg");
    $("iot-edit").classList.remove("hidden");
    $("iot-edit").scrollIntoView({ behavior: "smooth", block: "start" });
  }

  function closeIotForm() {
    $("iot-edit").classList.add("hidden");
  }

  function collectIotForm() {
    var idx = parseInt(val("iot-iot"), 10);
    var revived = val("iot-revivedata").split(",")
      .map(function (s) { return parseInt(s.trim(), 10); })
      .filter(function (n) { return !isNaN(n); });
    return {
      iot: idx,
      id: val("iot-id"), type: val("iot-type"), cron: val("iot-cron"),
      com: val("iot-com"), baudrate: parseInt(val("iot-baudrate"), 10) || 0,
      databits: parseInt(val("iot-databits"), 10) || 0,
      stopbits: parseInt(val("iot-stopbits"), 10) || 0,
      paritymode: parseInt(val("iot-paritymode"), 10) || 0,
      read_buff: val("iot-readbuff"), read_len: parseInt(val("iot-readlen"), 10) || 0,
      read_delay: parseInt(val("iot-readdelay"), 10) || 0,
      recive_len: parseInt(val("iot-reciven"), 10) || 0,
      revive_data_len: parseInt(val("iot-revivedatalen"), 10) || 0,
      revive_data: revived,
      offset: $("iot-offset").value === "true",
      ratio: parseFloat(val("iot-ratio")) || 0,
      crc: val("iot-crc")
    };
  }

  function saveIot() {
    var d = collectIotForm();
    if (isNaN(d.iot)) { showMsg("iot-edit-msg", "请输入传感器编号", "error"); return; }
    var req = iotEditMode === "edit" ? API.editIot(d.iot, d) : API.addIot(d.iot, d);
    req.then(function () {
      showMsg("iot-edit-msg", iotEditMode === "edit" ? "修改成功" : "添加成功", "success");
      loadIotList();
    }).catch(function (e) { showMsg("iot-edit-msg", e.message, "error"); });
  }

  function testIotRead() {
    var d = collectIotForm();
    if (isNaN(d.iot)) { showMsg("iot-edit-msg", "请输入传感器编号", "error"); return; }
    API.testRead(d).then(function (r) {
      var info = r.data;
      setVal("iot-readbuff", info.rawdata || "");
      showMsg("iot-edit-msg", "读取成功 原始数据: " + (info.rawdata || ""), "success");
    }).catch(function (e) { showMsg("iot-edit-msg", "读取失败: " + e.message, "error"); });
  }

  /* ================= 电源控制器 ================= */
  var powerHtmlCache = {};

  function renderPower(title, data) {
    var html = '<h3>' + esc(title) + '</h3>';
    if (data) {
      html += '<table class="table"><tbody>';
      Object.keys(data).forEach(function (k) {
        var v = data[k];
        if (Array.isArray(v)) v = v.join(", ");
        html += "<tr><td class='kv-label' style='width:220px'>" + esc(k) + "</td><td>" + esc(v) + "</td></tr>";
      });
      html += "</tbody></table>";
    }
    $("power-body").innerHTML = html;
    $("power-body").classList.remove("hidden");
  }

  function loadPower(para, title) {
    API.getPower(para).then(function (d) { renderPower(title, d.data); })
      .catch(function (e) { renderPower(title + "（读取失败）", null); showGlobal(e.message, "error"); });
  }

  function setPowerCapacityForm() {
    var type = parseInt(prompt("电池类型 (0=自定义,1=开口,2=密封,3=胶体,4=锂电):", "3"), 10);
    if (isNaN(type) || type < 0 || type > 4) { showGlobal("电池类型需为 0-4", "error"); return; }
    var cap = parseInt(prompt("电池容量 (Ah):", "100"), 10);
    if (isNaN(cap) || cap <= 0) { showGlobal("电池容量无效", "error"); return; }
    API.setPowerCapacity({ battertype: type, battercapacity: cap }).then(function () {
      showGlobal("电池类型/容量设置成功", "success");
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }

  /* ================= 历史数据 ================= */
  function fillHisIotSelect() {
    return API.getIots().then(function (d) {
      var sel = $("his-iot");
      sel.innerHTML = "";
      (d.data || []).forEach(function (iot) {
        var o = document.createElement("option");
        o.value = iot.iot;
        o.textContent = "#" + iot.iot + " " + iot.id;
        sel.appendChild(o);
      });
    });
  }
  function loadHis() {
    fillHisIotSelect().catch(function (e) { showGlobal(e.message, "error"); });
  }
  function queryHisDay() {
    var iot = parseInt(val("his-iot"), 10);
    var day = val("his-day");
    if (isNaN(iot) || !day) { showMsg("his-msg", "请选择传感器和日期", "error"); return; }
    API.getHisDay(iot, day).then(function (d) {
      renderHisTable(d.data || []);
    }).catch(function (e) { showMsg("his-msg", e.message, "error"); });
  }
  function queryHisOne() {
    var iot = parseInt(val("his-iot"), 10);
    var time = datetimeLocalToStr(val("his-time"));
    if (isNaN(iot) || !time) { showMsg("his-msg", "请选择传感器和时间", "error"); return; }
    API.getHisOne(iot, time).then(function (d) {
      renderHisTable(d.data ? [d.data] : []);
    }).catch(function (e) { showMsg("his-msg", e.message, "error"); });
  }
  function renderHisTable(list) {
    var body = $("his-table").querySelector("tbody");
    if (!list || list.length === 0) {
      $("his-table").classList.add("hidden");
      showMsg("his-msg", "未查询到数据", "info");
      return;
    }
    hideMsg("his-msg");
    var rows = "";
    list.forEach(function (r) {
      rows += "<tr><td>" + esc(r.time) + "</td><td>" + esc(r.id) + "</td><td>" + esc(r.type) +
        "</td><td>" + esc(r.value) + "</td><td>" + esc(r.offset) + "</td><td class='mono'>" + esc(r.rawdata) + "</td></tr>";
    });
    body.innerHTML = rows;
    $("his-table").classList.remove("hidden");
  }

  /* ================= 测试 ================= */
  function testSerial() {
    API.getIots().then(function (d) {
      var list = (d.data || []).filter(function (i) { return i.id === val("test-read-id"); });
      if (!list.length) { showGlobal("未找到该传感器配置", "error"); return; }
      API.testRead(list[0]).then(function (r) {
        var html = "";
        var info = r.data || {};
        ["id", "type", "value", "offset", "rawdata", "time"].forEach(function (k) {
          html += kvItem(k, info[k]);
        });
        $("test-read-body").innerHTML = html;
        $("test-read-body").classList.remove("hidden");
      }).catch(function (e) { showGlobal("串口读取失败: " + e.message, "error"); });
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }

  function bindEvents() {
    $("btn-login").addEventListener("click", doLogin);
    $("login-pass").addEventListener("keydown", function (e) { if (e.key === "Enter") doLogin(); });
    $("login-user").addEventListener("keydown", function (e) { if (e.key === "Enter") doLogin(); });
    $("btn-logout").addEventListener("click", doLogout);

    document.querySelectorAll("#nav .nav-item").forEach(function (a) {
      a.addEventListener("click", function (e) {
        e.preventDefault();
        location.hash = "#/" + a.dataset.page;
        navigate(a.dataset.page);
      });
    });

    window.addEventListener("hashchange", function () {
      var page = (location.hash || "").replace("#/", "") || "info";
      if (["info", "status", "base", "mqtt", "webapi", "iot", "power", "his", "test"].indexOf(page) >= 0) {
        navigate(page);
      }
    });

    $("info-refresh").addEventListener("click", loadInfo);
    $("status-refresh").addEventListener("click", loadStatus);
    $("base-save").addEventListener("click", function () {
      saveBase().then(function () { showGlobal("基础设置保存成功", "success"); }).catch(function (e) { showGlobal(e.message, "error"); });
    });
    $("mqtt-save").addEventListener("click", function () {
      saveMqtt().then(function () { showGlobal("MQTT 设置保存成功", "success"); }).catch(function (e) { showGlobal(e.message, "error"); });
    });
    $("webapi-save").addEventListener("click", function () {
      saveWebapi().then(function () { showGlobal("WebAPI 设置保存成功", "success"); }).catch(function (e) { showGlobal(e.message, "error"); });
    });

    $("iot-refresh").addEventListener("click", loadIotList);
    $("iot-add").addEventListener("click", function () { openIotForm("add", {}); });
    $("iot-load-type").addEventListener("click", loadTypeTemplate);
    $("iot-edit-save").addEventListener("click", saveIot);
    $("iot-edit-test").addEventListener("click", testIotRead);
    $("iot-edit-cancel").addEventListener("click", closeIotForm);

    $("power-info").addEventListener("click", function () { loadPower("info", "控制器设备信息"); });
    $("power-data").addEventListener("click", function () { loadPower("data", "控制器运行数据"); });
    $("power-set").addEventListener("click", function () { loadPower("set", "控制器参数"); });
    $("power-capacity").addEventListener("click", setPowerCapacityForm);

    $("his-day-query").addEventListener("click", queryHisDay);
    $("his-one-query").addEventListener("click", queryHisOne);

    $("test-mqtt").addEventListener("click", function () {
      API.getMqtt().then(function (d) {
        return API.testMqtt(d.data || {});
      }).then(function () { showGlobal("MQTT 测试成功", "success"); })
      .catch(function (e) { showGlobal(e.message, "error"); });
    });
    $("test-webapi").addEventListener("click", function () {
      API.getWebapi().then(function (d) {
        return API.testWebapi(d.data || {});
      }).then(function () { showGlobal("WebAPI 测试成功", "success"); })
      .catch(function (e) { showGlobal(e.message, "error"); });
    });
    $("test-read").addEventListener("click", testSerial);

    /* 测试页的传感器下拉框 */
    $("test-read-id").addEventListener("focus", function () {
      API.getIots().then(function (d) {
        var sel = $("test-read-id");
        if (sel.options.length) return;
        (d.data || []).forEach(function (iot) {
          var o = document.createElement("option");
          o.value = iot.id;
          o.textContent = "#" + iot.iot + " " + iot.id;
          sel.appendChild(o);
        });
      }).catch(function () {});
    });

    window.App = window.App || {};
    App.editIotFromList = function (n) {
      API.getIot(n).then(function (d) { openIotForm("edit", d.data || {}); })
        .catch(function (e) { showGlobal(e.message, "error"); });
    };
    App.delIotFromList = function (n) {
      if (!confirm("确定删除传感器 #" + n + " 吗？")) return;
      API.delIot(n).then(function () {
        showGlobal("删除成功", "success");
        loadIotList();
      }).catch(function (e) { showGlobal(e.message, "error"); });
    };
  }

  /* ---------- 传感器类型模板 (help) ---------- */
  function loadTypeTemplate() {
    hideMsg("iot-edit-msg");
    API.getHelpType().then(function (d) {
      var types = d.data || [];
      if (!types.length) { showGlobal("未找到类型模板", "info"); return; }
      var lines = types.map(function (t, i) {
        return (i + 1) + ". [" + t.OrgName + " / " + t.Name + "] type=" + t.Type + " 文件=" + t.File;
      });
      var pick = prompt("选择设备类型编号:\n\n" + lines.join("\n"));
      var idx = parseInt(pick, 10) - 1;
      if (isNaN(idx) || !types[idx]) { showGlobal("编号无效", "error"); return; }
      API.getHelpRead(types[idx].File).then(function (r) {
        var template = r.data || {};
        iotNextIndex().then(function (nextId) {
          template.iot = nextId;
          openIotForm("add", template);
          showMsg("iot-edit-msg", "已从模板加载（文件: " + types[idx].File + "），请确认编号与参数后保存", "info");
        });
      }).catch(function (e) { showGlobal(e.message, "error"); });
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }

  function iotNextIndex() {
    /* 返回当前最大的传感器编号 + 1 */
    return API.getIots().then(function (d) {
      var maxId = (d.data || []).reduce(function (m, iot) {
        return (typeof iot.iot === "number" && iot.iot > m) ? iot.iot : m;
      }, 0);
      return maxId + 1;
    });
  }

  function bootstrap() {
    bindEvents();
    var page = (location.hash || "").replace("#/", "") || "info";
    if (["info", "status", "base", "mqtt", "webapi", "iot", "power", "his", "test"].indexOf(page) >= 0) {
      enterMain();
      navigate(page);
    } else {
      $("login-page").classList.remove("hidden");
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", bootstrap);
  } else {
    bootstrap();
  }
})();