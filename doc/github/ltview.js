if (!Array.prototype.find) {
  Array.prototype.find = function(predicate) {
     // 1. Let O be ? ToObject(this value).
      if (this == null) {
        throw new TypeError('"this" is null or not defined');
      }

      var o = Object(this);

      // 2. Let len be ? ToLength(? Get(O, "length")).
      var len = o.length >>> 0;

      // 3. If IsCallable(predicate) is false, throw a TypeError exception.
      if (typeof predicate !== 'function') {
        throw new TypeError('predicate must be a function');
      }

      // 4. If thisArg was supplied, let T be thisArg; else let T be undefined.
      var thisArg = arguments[1];

      // 5. Let k be 0.
      var k = 0;

      // 6. Repeat, while k < len
      while (k < len) {
        // a. Let Pk be ! ToString(k).
        // b. Let kValue be ? Get(O, Pk).
        // c. Let testResult be ToBoolean(? Call(predicate, T, « kValue, k, O »)).
        // d. If testResult is true, return kValue.
        var kValue = o[k];
        if (predicate.call(thisArg, kValue, k, o)) {
          return kValue;
        }
        // e. Increase k by 1.
        k++;
      }

      // 7. Return undefined.
      return undefined;
    }
}
if (!Array.prototype.forEach) {
  Array.prototype.forEach = function(callback/*, thisArg*/) {

    var T, k;

    if (this == null) {
      throw new TypeError('this is null or not defined');
    }

    // 1. Let O be the result of calling toObject() passing the
    // |this| value as the argument.
    var O = Object(this);

    // 2. Let lenValue be the result of calling the Get() internal
    // method of O with the argument "length".
    // 3. Let len be toUint32(lenValue).
    var len = O.length >>> 0;

    // 4. If isCallable(callback) is false, throw a TypeError exception. 
    // See: http://es5.github.com/#x9.11
    if (typeof callback !== 'function') {
      throw new TypeError(callback + ' is not a function');
    }

    // 5. If thisArg was supplied, let T be thisArg; else let
    // T be undefined.
    if (arguments.length > 1) {
      T = arguments[1];
    }

    // 6. Let k be 0.
    k = 0;

    // 7. Repeat while k < len.
    while (k < len) {

      var kValue;

      // a. Let Pk be ToString(k).
      //    This is implicit for LHS operands of the in operator.
      // b. Let kPresent be the result of calling the HasProperty
      //    internal method of O with argument Pk.
      //    This step can be combined with c.
      // c. If kPresent is true, then
      if (k in O) {

        // i. Let kValue be the result of calling the Get internal
        // method of O with argument Pk.
        kValue = O[k];

        // ii. Call the Call internal method of callback with T as
        // the this value and argument list containing kValue, k, and O.
        callback.call(T, kValue, k, O);
      }
      // d. Increase k by 1.
      k++;
    }
    // 8. return undefined.
  };
}
function toArray(list) {
    var _list = [];
    for( var i = 0; i < list.length; i++) {
        _list.push(list[i])
    }
    return _list
}
window.onload = function (){
    var list = toArray(document.querySelectorAll('.md-toc-item')), oldMax, oldname = 0;
    var btn = document.querySelector('.toggleExpand');
    var isMobile = window.innerWidth < 768;
    if (isMobile) {
        var menu = this.document.querySelector('.md-toc');
        if (!btn) {
            btn = document.createElement('button')
            btn.setAttribute('open', 'false')
            btn.className = 'toggleExpand'
            btn.innerText = '>'
            menu.appendChild(btn)
        }
        addStyle('.toggleExpand', {
            display: 'block',
            position: 'absolute',
            right: 0,
            top: 0,
            width: '30px',
            height: '50px',
            'border': '1px solid #ccc',
            'border-right-width': '0',
            'background': '#fafafa',
            'font-size': '22px',
            'font-family': 'cursive',
            'font-weight': '600',
            'outline': 'none',
            color: 'rgb(65, 131, 196)'
        })
        document.querySelector('#write').onclick = function(e) {
            var btn = document.querySelector('.toggleExpand')
            var iso = btn.getAttribute('open') == 'true'
            var ismenu = e.target.className.indexOf('md-toc-item') != -1 || e.target.className.indexOf('md-toc-inner') != -1
            if(iso && e.target.className != 'toggleExpand' && !ismenu) {
                btn.setAttribute('open', false)
                btn.innerText = '>'
                addStyle('.md-toc', {
                    'left': '-230px',
                    'overflow': 'hidden'
                })
            }
            if (e.target.className == 'toggleExpand') {
                if(iso) {
                    btn.setAttribute('open', false)
                    btn.innerText = '>'
                    addStyle('.md-toc', {
                        'left': '-230px',
                        'overflow': 'hidden'
                    })
                } else {
                    btn.setAttribute('open', true)
                    btn.innerText = '<'
                    addStyle('.md-toc', {
                        'left': '20px',
                        'overflow': 'auto'
                    })
                }
            }
        }
    } else {
        btn && (btn.style.display = 'none')
    }

    addStyle('#write', {margin: 0, 'word-break': 'keep-all'})
    addStyle('.typora-export', {
        'margin-left': isMobile ? 0 :'300px'
    })
    addStyle('.md-toc', {
        'position': 'fixed',
        'background-color': '#ffffff',
        'height': '100%',
        'left': isMobile ? '-230px' : '20px',
        'top': '0',
        'margin-top': '0px',
        'border-right': '1px solid #ccc',
        'border-radius': 0,
        'overflow': isMobile ? 'hidden' : 'auto',
        'z-index': '100',
        'width': '260px',
        'padding-right': '30px',
        'transition': 'left .3s'
    })
    addStyle('.md-toc a', {
        'display': 'block',
        'min-width': '240px',
        'word-break': 'keep-all'
    })
    list.forEach(function (listItem, listInd){
        listItem.className != 'md-toc-item md-toc-h1' && (listItem.style.display = 'none')
        listItem.onclick = function (){
            getList(listItem.className, listInd)
        }
    })
    function getList(name, ind){
        var showName = name.substr(0, name.length-1) + (parseInt(name.substr(-1)) + 1)
        var indList = [], max;
        list.forEach(function (listItem, listInd){
            listItem.className == name && indList.push(listInd)
        })
        max = ind == indList[indList.length-1] ? list.length : indList.find(function(item){
            return item > ind
        })
        num = name.substr(-1);
        list.forEach(function (listItem, listInd){
            if(listInd > ind && listInd < max){
                if (listItem.className.substr(-1) >= (parseInt(num)+1) && (listItem.style.display === 'block')) {
                    listItem.style.display = 'none';
                } else {
                    listItem.className == showName && (listItem.style.display = 'block')
                }
            }
        })
    }
    function addStyle (selector, styles){
        var ele = document.querySelectorAll(selector)
        for(var i = 0; i < ele.length; i++) {
            for (var key in styles) {
                if (styles.hasOwnProperty(key)) {
                    try{
                        ele[i].style[key] = styles[key]
                    } catch(e) {
                        console.log(key, e);
                    }
                }
            }
        }
    }
}
window.onresize = window.onload