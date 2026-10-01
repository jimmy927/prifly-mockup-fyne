// Starts the Go wasm program from the bytes the packer embedded (no fetch, no network).
(function () {
  // Fyne 2.8 opens an IndexedDB at start and panics on a nil request when indexedDB.open throws,
  // which it does in a file/opaque-origin page. Make open() fail the polite way instead.
  var usable = false;
  try { usable = !!window.indexedDB && !!window.indexedDB.open('mockup-probe'); } catch (e) { usable = false; }
  if (!usable) {
    var fake = {
      open: function () {
        // an object that passes `instanceof IDBRequest` (go-indexeddb checks) but owns every member
        var req = Object.create(IDBOpenDBRequest.prototype);
        var listeners = {};
        Object.defineProperties(req, {
          error: { value: new DOMException('IndexedDB is not available in this page', 'InvalidStateError') },
          readyState: { value: 'done' },
          result: { value: null },
          addEventListener: { value: function (type, fn) { (listeners[type] = listeners[type] || []).push(fn); } },
          removeEventListener: { value: function (type, fn) { listeners[type] = (listeners[type] || []).filter(function (f) { return f !== fn; }); } }
        });
        setTimeout(function () {
          (listeners.error || []).slice().forEach(function (fn) { fn({ type: 'error', target: req, preventDefault: function () {}, stopPropagation: function () {} }); });
        }, 0);
        return req;
      }
    };
    Object.defineProperty(window, 'indexedDB', { value: fake, configurable: true });
  }

  // Mockup control, outside the app frame: a 34 px strip on top, the Fyne canvas below it.
  // Fyne sizes its canvas buffer once from window.innerHeight, so that is told 766 (= 800 - 34) only until
  // the canvas exists, then given back: the mockup tool sizes its picture from innerHeight and would crop it.
  var BAR = 34, APP_H = 766;
  var realInner = Object.getOwnPropertyDescriptor(window, 'innerHeight');
  Object.defineProperty(window, 'innerHeight', { configurable: true, get: function () { return APP_H; } });
  var restore = setInterval(function () {
    if (!document.querySelector('canvas')) return;
    clearInterval(restore);
    Object.defineProperty(window, 'innerHeight', realInner);
  }, 5);
  var css = document.createElement('style');
  css.textContent = 'html,body{box-sizing:border-box!important}body{padding-top:' + BAR + 'px!important}canvas{display:block!important;height:' + APP_H + 'px!important}' +
    '#mockup-control{position:fixed;top:0;left:0;right:0;height:' + BAR + 'px;box-sizing:border-box;display:flex;align-items:center;gap:6px;' +
    'padding:0 10px;background:#1f2430;color:#cfd6e4;font:12px system-ui,sans-serif;z-index:2147483001}' +
    '#mockup-control a{padding:3px 10px;border-radius:12px;color:#fff;text-decoration:none;background:#3a4256}' +
    '#mockup-control a.on{background:#2f7bf0}#mockup-control b{font-weight:600;margin-right:6px}';
  document.head.appendChild(css);
  var bar = document.createElement('div');
  bar.id = 'mockup-control';
  bar.innerHTML = '<b>Mockup control (not part of the app)</b>state:';
  [['Notes', 'notes'], ['Editor', 'editor'], ['Settings', 'settings'], ['About', 'about'], ['Delete dialog', 'dialog']].forEach(function (t) {
    var a = document.createElement('a');
    a.textContent = t[0]; a.href = '#' + t[1];
    if ((window.__state || 'notes') === t[1]) a.className = 'on';
    bar.appendChild(a);
  });
  document.body.appendChild(bar);

  // The picture waits for Fyne's own first frame: the Go side calls __goReady 1.5 s after the
  // window is shown; 7.5 s at the latest.
  window.__mockupReady = new Promise(function (resolve) {
    window.__goReady = function () { resolve(true); };
    setTimeout(function () { resolve(true); }, 7500);
  });
  var go = new Go();
  window.__mockupBytes('Quillnote.wasm')
    .then(function (bytes) { return WebAssembly.instantiate(bytes, go.importObject); })
    .then(function (result) { return go.run(result.instance); })
    .catch(function (e) { document.body.textContent = 'Fyne failed to start: ' + e; });
})();
