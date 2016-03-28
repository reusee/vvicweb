import {div, input, label, button, none, span, img,
  e, Store, Component} from './base'
import {call} from './call'

let initState = {
};

class App extends Component {
  render(state) {
    return div({
      style: {
        'font-size': '14px',
      },
    }, [
      e(Control),
      !state.info ? none : e(Info, state.info),
    ]);
  }
}

class Control extends Component {
  render(state) {
    return div({}, [
      label({}, [
        '货号',
        input({
          type: 'text',
          id: 'good-number',
          onclick: () => {
            document.getElementById('good-number').select();
          },
        }),
      ]),
      button({
        onclick: () => {
          let id = parseInt(document.getElementById('good-number').value);
          if (id) {
            emit(ev_get_info, id);
          }
        },
      }, '提交'),
    ]);
  }
}

let section_style = {
  margin: '20px 0',
}

let image_style = {
  'max-width': '150px',
  margin: '10px',
  padding: '10px',
  border: '1px solid #CCC',
}

class Info extends Component {
  render(state) {
    let price = parseFloat(state.Price);
    return div({}, [
      // 标题等
      div({
        style: {...section_style,
        },
      }, [
        div({}, '梦丹铃 2016' + state.Title),
        div({}, '拿货价格：￥' + state.Price),
        div({}, [
          '最高成本：￥',
          price + 18 + (price * 0.08) + 13 + (price * 0.002),
        ]),
        div({}, '货号：' + state.Id),
      ]),
      // 下载
      div({
        style: {...section_style,
        },
      }, [
        e('a', {
          href: 'http://127.0.0.1:7899/download/' + state.Id,
        }, '下载图包'),
      ]),
      // 属性
      div({
        style: {...section_style
        },
      }, state.Attrs.map((entry) => {
        return div({
          style: {
            width: '30%',
            clear: 'both',
          },
        }, [
          span({
            style: {
              float: 'left',
            },
          }, entry[0]),
          span({
            style: {
              float: 'right',
            },
          }, entry[1]),
        ]);
      })),
      // 商品图
      div({
        style: {...section_style,
        },
      }, state.ThemeImages.map((src) => {
        return img({
          src: src,
          style: {...image_style,
          },
        });
      })),
      // 详情图
      div({
        style: {...section_style,
        },
      }, state.DetailImages.map((src) => {
        return img({
          src: src,
          style: {...image_style,
          },
        });
      })),
    ]);
  }
}

let app = new App(initState);
app.bind(document.getElementById('app'));

let store = new Store(initState);
store.setComponent(app);

export let emit = store.emit.bind(store);

export let ev_get_info = (state, id) => {
  call('GetInfo', {
    id: id,
  }, (resp) => {
    emit(ev_update_info, resp);
  });
};

export let ev_update_info = (state, info) => {
  console.log(info);
  if (!info.ok) {
    return;
  }
  return {...state,
    info: info,
  };
};
