/* RTU 设备管理 - 接口封装层
 *
 * 接口约定：
 *   rtucgi.cgi 为 CGI 程序，必须通过 boa 的 ScriptAlias 目录访问，
 *   本设备为  /usr/boa/www/cgi-bin/rtucgi.cgi，浏览器地址为 /cgi-bin/rtucgi.cgi。
 *   除 login 外均需在查询串携带 token 参数：?action=xxx&token=<jwt>。
 *   成功返回 {code:200,...}，失败返回 {code:400,msg}。
 *
 * CGI 地址可通过以下方式指定（优先级从高到低）：
 *   1. index.html 的 ?cgi= 查询参数，如 index.html?cgi=/custom/rtucgi.cgi
 *   2. 本文件下方 CONFIG.cgiUrl 手动填写
 *   3. 默认 /cgi-bin/rtucgi.cgi
 *
 * 注意：这里使用 XMLHttpRequest 而非 fetch，因为目标设备上的 boa 0.94
 *   对部分浏览器的 fetch(POST, json) 请求会返回 400 Bad Request，
 *   XHR 显式设置 Content-Type 并携带正确 Content-Length，兼容性最好。
 */
var CONFIG = { cgiUrl: "" };

(function () {
  var override = location.search.match(/[?&]cgi=([^&]+)/);
  if (override) {
    CONFIG.cgiUrl = decodeURIComponent(override[1]);
  } else {
    CONFIG.cgiUrl = "/cgi-bin/rtucgi.cgi";
  }
})();

var API = (function () {
  var tokenKey = "rtu_token";

  function getToken() { return localStorage.getItem(tokenKey) || ""; }
  function setToken(t) { localStorage.setItem(tokenKey, t); }
  function clearToken() { localStorage.removeItem(tokenKey); }

  function buildUrl(action, params) {
    var url = CONFIG.cgiUrl + "?action=" + encodeURIComponent(action);
    if (params) {
      Object.keys(params).forEach(function (k) {
        url += "&" + encodeURIComponent(k) + "=" + encodeURIComponent(params[k]);
      });
    }
    return url;
  }

  /* 统一的 CGI 请求。body 为对象时自动 JSON 序列化。
   * 除 login 外自动附加 token 查询参数（CGI 从查询串校验认证）。 */
  function request(action, params, method, body) {
    var query = {};
    for (var k in (params || {})) query[k] = params[k];
    if (action !== "login") {
      query.token = getToken();
    }
    var url = buildUrl(action, query);
    var hasBody = body !== undefined && body !== null;
    return new Promise(function (resolve, reject) {
      var xhr = new XMLHttpRequest();
      xhr.open(method || "GET", url, true);
      if (hasBody) {
        xhr.setRequestHeader("Content-Type", "application/json");
      }
      xhr.setRequestHeader("Accept", "application/json");

      xhr.onreadystatechange = function () {
        if (xhr.readyState !== 4) return;
        var text = xhr.responseText || "";
        var data = null;
        try { data = JSON.parse(text); } catch (e) { data = text; }

        if (xhr.status !== 200) {
          var msg;
          if (typeof data === "object" && data !== null) {
            msg = data.msg || ("HTTP " + xhr.status);
          } else {
            msg = text || ("HTTP " + xhr.status);
          }
          var err = new Error(msg);
          err.status = xhr.status;
          err.url = url;
          err.data = data;
          reject(err);
          return;
        }
        /* 服务端把 HTTP 状态码写入 JSON 的 code 字段 */
        if (data && typeof data === "object" && data.code !== 200) {
          var e2 = new Error(data.msg || "操作失败");
          e2.status = xhr.status;
          e2.data = data;
          reject(e2);
          return;
        }
        resolve(data);
      };
      xhr.onerror = function () {
        reject(new Error("网络错误，无法连接 " + url));
      };

      xhr.send(hasBody ? JSON.stringify(body) : null);
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
    getToken: getToken, clearToken: clearToken,
    cgiUrl: CONFIG.cgiUrl
  };
})();