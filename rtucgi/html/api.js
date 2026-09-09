/* RTU 设备管理 - 接口封装层
 * 所有接口均为 CGI 查询串调用：rtucgi.cgi?action=xxx&token=<jwt>
 * 除 login 外均需携带 token。成功时返回 {code:200,...}，失败返回 {code:400,msg}。
 * 可通过下方 CONFIG.cgiUrl 修改后端 CGI 地址。
 */
var CONFIG = {
  cgiUrl: "rtucgi.cgi"
};

var API = (function () {
  var tokenKey = "rtu_token";

  function getToken() { return localStorage.getItem(tokenKey) || ""; }
  function setToken(t) { localStorage.setItem(tokenKey, t); }
  function clearToken() { localStorage.removeItem(tokenKey); }

  function buildUrl(action, params) {
    var url = CONFIG.cgiUrl + "?action=" + encodeURIComponent(action);
    if (params) {
      for (var k in params) {
        url += "&" + encodeURIComponent(k) + "=" + encodeURIComponent(params[k]);
      }
    }
    return url;
  }

  function request(action, params, method, body) {
    var url = buildUrl(action, params);
    var opt = { method: method || "GET" };
    if (body !== undefined) {
      opt.headers = { "Content-Type": "application/json" };
      opt.body = JSON.stringify(body);
    }
    return fetch(url, opt).then(function (resp) {
      return resp.text().then(function (text) {
        var data = null;
        try { data = JSON.parse(text); } catch (e) { data = text; }
        if (!resp.ok) {
          if (typeof data === "object" && data !== null) {
            var err = new Error(data.msg || "请求失败");
            err.data = data;
            throw err;
          }
          throw new Error(data || "HTTP " + resp.status);
        }
        return data;
      }).then(function (data) {
        /* 服务端又把 HTTP 状态码写入 JSON 的 code 字段 */
        if (data && typeof data === "object" && data.code !== 200) {
          var err = new Error(data.msg || "操作失败");
          err.data = data;
          throw err;
        }
        return data;
      });
    });
  }

  /* ---------- 登录 ---------- */
  function login(username, password) {
    return request("login", null, "POST", { username: username, password: password })
      .then(function (d) { setToken(d.data); return d.data; });
  }

  /* ---------- 基础设置 base ---------- */
  function getBase()   { return request("base"); }
  function saveBase(d) { return request("base", null, "POST", d); }

  /* ---------- MQTT ---------- */
  function getMqtt()   { return request("mqtt"); }
  function saveMqtt(d) { return request("mqtt", null, "POST", d); }

  /* ---------- WebAPI ---------- */
  function getWebapi()   { return request("webapi"); }
  function saveWebapi(d) { return request("webapi", null, "POST", d); }

  /* ---------- 传感器 iot ---------- */
  function getIots()          { return request("iot"); }
  function getIot(n)          { return request("iot", { iot: n }); }
  function addIot(n, d)       { return request("iot", { iot: n, para: "add" }, "POST", d); }
  function editIot(n, d)      { return request("iot", { iot: n, para: "edit" }, "POST", d); }
  function delIot(n)          { return request("iot", { iot: n, para: "del" }, "POST"); }

  /* ---------- 运行状态 / 系统信息 ---------- */
  function getStatus()        { return request("status"); }
  function getInfo()          { return request("info"); }

  /* ---------- 电源控制器 power ---------- */
  function getPower(para)     { return request("power", { para: para }); }
  function setPowerCapacity(d){ return request("power", null, "POST", d); }

  /* ---------- 历史数据 his ---------- */
  function getHisDay(iot, time){ return request("his", { para: "day", iot: iot, time: time }); }
  function getHisOne(iot, time){ return request("his", { para: "one", iot: iot, time: time }); }

  /* ---------- 测试 test ---------- */
  function testRead(cfg)      { return request("test", { para: "read" }, "POST", cfg); }
  function testMqtt(cfg)      { return request("test", { para: "mqtt" }, "POST", cfg); }
  function testWebapi(cfg)    { return request("test", { para: "webapi" }, "POST", cfg); }

  /* ---------- 帮助/模板 help ---------- */
  function getHelpOrg()       { return request("help", { para: "org" }); }
  function getHelpType()      { return request("help", { para: "type" }); }
  function getHelpRead(file)  { return request("help", { para: "read", file: file }); }
  function getHelpMqtt()      { return request("help", { para: "mqtt" }); }
  function getHelpWebapi()    { return request("help", { para: "webapi" }); }

  return {
    login: login,
    getBase: getBase, saveBase: saveBase,
    getMqtt: getMqtt, saveMqtt: saveMqtt,
    getWebapi: getWebapi, saveWebapi: saveWebapi,
    getIots: getIots, getIot: getIot,
    addIot: addIot, editIot: editIot, delIot: delIot,
    getStatus: getStatus, getInfo: getInfo,
    getPower: getPower, setPowerCapacity: setPowerCapacity,
    getHisDay: getHisDay, getHisOne: getHisOne,
    testRead: testRead, testMqtt: testMqtt, testWebapi: testWebapi,
    getHelpOrg: getHelpOrg, getHelpType: getHelpType, getHelpRead: getHelpRead,
    getHelpMqtt: getHelpMqtt, getHelpWebapi: getHelpWebapi,
    getToken: getToken, clearToken: clearToken
  };
})();