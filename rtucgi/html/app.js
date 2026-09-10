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

  function esc(s) {
    return String(s == null ? "" : s).replace(/[&<>"']/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c];
    });
  }
  function kvItem(label, value) {
    return '<div class="kv-item"><div class="kv-label">' + esc(label) +
      '</div><div class="kv-value">' + esc(value) + "</div></div>";
  }
  function fmtBool(b) { return b ? "是" : "否"; }

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
  var pageTimers = { info: null, status: null, iot: null };
  function clearTimers() {
    Object.keys(pageTimers).forEach(function (k) {
      if (pageTimers[k]) { clearInterval(pageTimers[k]); pageTimers[k] = null; }
    });
  }
  function scheduleAutoRefresh(page, fn, ms) {
    if (pageTimers[page]) clearInterval(pageTimers[page]);
    pageTimers[page] = setInterval(fn, ms);
  }
  function navigate(page) {
    currentPage = page;
    clearTimers();
    document.querySelectorAll("#nav .nav-item").forEach(function (a) {
      a.classList.toggle("active", a.dataset.page === page);
    });
    document.querySelectorAll("#nav .nav-group").forEach(function (g) {
      var inGroup = (g.dataset.pages || "").split(",").indexOf(page) >= 0;
      g.querySelector(".nav-toggle").classList.toggle("active", inGroup);
      g.querySelector(".nav-toggle").classList.toggle("open", inGroup);
      g.querySelector(".nav-sub").classList.toggle("hidden", !inGroup);
    });
    document.querySelectorAll(".content .page").forEach(function (p) {
      p.classList.toggle("hidden", p.id !== "page-" + page);
    });
    loadPage(page);
  }

  function toggleNavGroup(toggle) {
    var sub = document.getElementById(toggle.dataset.toggle);
    if (!sub) return;
    sub.classList.toggle("hidden");
    toggle.classList.toggle("open", !sub.classList.contains("hidden"));
  }

  function loadPage(page) {
    switch (page) {
      case "info": loadInfo(); scheduleAutoRefresh("info", loadInfo, 60000); break;
      case "status": loadStatus(); scheduleAutoRefresh("status", loadStatus, 30000); break;
      case "base": loadBase(); break;
      case "network": loadNetwork(); break;
      case "mqtt": loadMqtt(); break;
      case "webapi": loadWebapi(); break;
      case "iot": loadIotList(); scheduleAutoRefresh("iot", loadIotList, 60000); break;
      case "power": break; /* 由按钮触发 */
      case "his": loadHis(); break;
      case "test": break;
    }
  }

  /* ================= 格式化工具 ================= */
  function fmtBytes(v) {
    v = Number(v) || 0;
    if (v >= 1073741824) return (v / 1073741824).toFixed(2) + " GB";
    if (v >= 1048576) return (v / 1048576).toFixed(1) + " MB";
    if (v >= 1024) return (v / 1024).toFixed(1) + " KB";
    return v + " B";
  }
  function fmtUptime(s) {
    s = Math.floor(Number(s) || 0);
    var d = Math.floor(s / 86400);
    var h = Math.floor((s % 86400) / 3600);
    var m = Math.floor((s % 3600) / 60);
    var str = "";
    if (d) str += d + "天";
    if (h || d) str += h + "小时";
    str += m + "分";
    return str;
  }
  function fmtDur(starttime) {
    if (!starttime) return "-";
    var t = new Date(String(starttime).replace(" ", "T"));
    if (isNaN(t.getTime())) return "-";
    return fmtUptime((Date.now() - t.getTime()) / 1000);
  }
  function valOrDash(v) {
    return (v === "" || v === null || v === undefined) ? "无" : v;
  }
  function nowStr() {
    function p(x) { return (x < 10 ? "0" : "") + x; }
    var d = new Date();
    return d.getFullYear() + "-" + p(d.getMonth() + 1) + "-" + p(d.getDate()) +
      " " + p(d.getHours()) + ":" + p(d.getMinutes()) + ":" + p(d.getSeconds());
  }

  /* ================= 系统信息 ================= */
  function loadInfo() {
    API.getInfo().then(function (d) {
      var s = d.data || {};
      function section(title, items) {
        var html = '<div class="card"><h3>' + esc(title) + '</h3><div class="kv-grid">';
        items.forEach(function (it) { html += kvItem(it[0], it[1]()); });
        html += '</div></div>';
        return html;
      }
      var html = "";
      function cleanIp(ip) { return valOrDash(String(ip || "").replace(/\/\d+$/, "")); }
      html += section("版本信息", [
        ["系统版本", function () { return valOrDash(s.version); }],
        ["采集程序", function () { return valOrDash(s.read_version); }],
        ["客户端程序", function () { return valOrDash(s.rtu_version); }],
        ["Web 管理", function () { return valOrDash(s.cgi_version); }]
      ]);
      html += section("运行与资源", [
        ["已运行时间", function () { return fmtUptime(s.uptime); }],
        ["内存总量", function () { return fmtBytes(s.memtotal); }],
        ["空闲内存", function () { return fmtBytes(s.memfree); }],
        ["内存使用率", function () {
          return s.memtotal ? (Math.round((1 - s.memfree / s.memtotal) * 1000) / 10) + "%" : "-";
        }],
        ["磁盘总量", function () { return fmtBytes(s.disktotal); }],
        ["磁盘剩余", function () { return fmtBytes(s.diskfree); }],
        ["存储使用率", function () {
          return s.disktotal ? (Math.round((1 - s.diskfree / s.disktotal) * 1000) / 10) + "%" : "-";
        }]
      ]);
      html += section("网络", [
        ["MAC 地址", function () { return valOrDash(s.mac); }],
        ["IP 地址", function () { return cleanIp(s.ip); }],
        ["子网掩码", function () { return valOrDash(s.mask); }],
        ["网关", function () { return valOrDash(s.gateway); }],
        ["4G 内网 IP", function () { return valOrDash(s.gprs_ip); }]
      ]);
      var has4g = s.gprs_imsi || s.gprs_ccid || s.gprs_csq;
      html += section("4G 通信", has4g ? [
        ["IMSI", function () { return valOrDash(s.gprs_imsi); }],
        ["CCID (SIM 卡号)", function () { return valOrDash(s.gprs_ccid); }],
        ["信号强度 CSQ", function () { return valOrDash(s.gprs_csq); }]
      ] : [["状态", function () { return "未接入 4G 模块"; }]]);
      $("info-body").innerHTML = html;
    }).catch(function (e) {
      $("info-body").innerHTML = '<div class="msg error">' + esc(e.message) + "</div>";
    });
  }

  /* ================= 运行状态 ================= */
  function loadStatus() {
    API.getStatus().then(function (d) {
      var s = d.data || {};
      var html = "";
      if (!s.runing) {
        html += '<div class="msg error">采集程序（iotread）未运行，无法正常采集数据。</div>';
      } else {
        html += '<div class="msg success">采集程序运行中</div>';
      }
      html += '<div class="card"><div class="kv-grid">';
      html += kvItem("启动时间", s.starttime);
      html += kvItem("已运行时长", fmtDur(s.starttime));
      html += kvItem("采集总次数", s.count);
      html += kvItem("读取成功次数", s.readcount);
      html += kvItem("MQTT 发送次数", s.mqttcount);
      html += kvItem("WebAPI 发送次数", s.webapicount);
      html += kvItem("程序版本", s.version);
      html += '</div></div>';
      $("status-body").innerHTML = html;
    }).catch(function (e) {
      $("status-body").innerHTML = '<div class="msg error">' + esc(e.message) + "</div>";
    });
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

  /* ================= 网络设置 ================= */
  function loadNetwork() {
    API.getNetwork().then(function (d) {
      var n = d.data || {};
      setVal("net-method", n.Method || n.method || "");
      setVal("net-mac", n.Mac || n.mac || "");
      setVal("net-ip", n.ip);
      setVal("net-mask", n.mask);
      setVal("net-gateway", n.gateway);
    }).catch(function (e) { showGlobal(e.message, "error"); });
  }
  function saveNetwork() {
    return API.saveNetwork({
      ip: val("net-ip"), mask: val("net-mask"), gateway: val("net-gateway")
    });
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
      var latests = list.map(function (iot) {
        return API.getHisOne(iot.iot, nowStr()).then(function (r) {
          iot._latest = r.data || null;
        }).catch(function () { iot._latest = null; });
      });
      return Promise.all(latests).then(function () {
        var rows = "";
        list.forEach(function (iot) {
          var latest = iot._latest;
          var hasData = latest && latest.time && String(latest.time).indexOf("0001-") !== 0 && latest.rawdata;
          var valCell = hasData
            ? latest.value + "<div class='small'>" + esc(latest.time) + "</div>"
            : '<span class="warn-text">无数据</span>';
          rows += "<tr>" +
            "<td>" + iot.iot + "</td>" +
            "<td>" + esc(iot.id) + "</td>" +
            "<td>" + esc(iot.type) + "</td>" +
            "<td>ttyS" + esc(iot.com) + "</td>" +
            "<td class='mono'>" + esc(iot.cron) + "</td>" +
            "<td>" + valCell + "</td>" +
            "<td class='mono small'>" + (hasData ? esc(latest.rawdata) : "") + "</td>" +
            "<td>" +
              "<button class='btn btn-secondary btn-sm' onclick='App.editIotFromList(" + iot.iot + ")'>编辑</button> " +
              "<button class='btn btn-danger btn-sm' onclick='App.delIotFromList(" + iot.iot + ")'>删除</button>" +
            "</td></tr>";
        });
        $("iot-list").querySelector("tbody").innerHTML = rows;
      });
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
      var info = r.data || {};
      showMsg("iot-edit-msg",
        "读取成功：数值 " + info.value +
        "，时间 " + (info.time || "-") +
        "，原始数据 " + (info.rawdata || "-"), "success");
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

    var navEl = $("nav");
    navEl.addEventListener("click", function (e) {
      var t = e.target, a = null;
      if (t && t.closest) {
        a = t.closest(".nav-item");
      } else {
        var n = t;
        while (n && n !== navEl && n.classList && !n.classList.contains("nav-item")) n = n.parentNode;
        if (n && n !== navEl && n.classList && n.classList.contains("nav-item")) a = n;
      }
      if (!a) return;
      e.preventDefault();
      if (a.classList.contains("nav-toggle")) {
        toggleNavGroup(a);
      } else if (a.dataset.page) {
        location.hash = "#/" + a.dataset.page;
        navigate(a.dataset.page);
      }
    });

    window.addEventListener("hashchange", function () {
      var page = (location.hash || "").replace("#/", "") || "info";
      if (["info", "status", "base", "network", "mqtt", "webapi", "iot", "power", "his", "test"].indexOf(page) >= 0) {
        navigate(page);
      }
    });

    $("info-refresh").addEventListener("click", loadInfo);
    $("status-refresh").addEventListener("click", loadStatus);
    $("base-save").addEventListener("click", function () {
      saveBase().then(function () { showGlobal("基础设置保存成功", "success"); }).catch(function (e) { showGlobal(e.message, "error"); });
    });
    $("net-save").addEventListener("click", function () {
      saveNetwork().then(function () { showGlobal("网络设置保存成功（重启生效）", "success"); }).catch(function (e) { showGlobal(e.message, "error"); });
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
    if (["info", "status", "base", "network", "mqtt", "webapi", "iot", "power", "his", "test"].indexOf(page) >= 0) {
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