import React, { Component } from 'react'
import {
  BrowserRouter as Router,
  Route
} from 'react-router-dom'

import './App.css'
import { Layout } from 'antd'
import LeftMenu from './components/LeftMenu'
import Label from './components/Label'
import Push from './components/Push'
import Pop from './components/Pop'
import Stock from './components/Stock'

const { Content, Footer, Sider } = Layout

class App extends Component {
  render () {
    return (
      <Router>
        <Layout>
          <Sider style={{ overflow: 'auto', height: '100vh', position: 'fixed', left: 0 }}>
            {/* <div className='logo' /> */}
            <LeftMenu />
          </Sider>
          <Layout style={{ marginLeft: 200 }}>
            <Content style={{ margin: '24px 16px 0', overflow: 'initial' }}>
              <div style={{ padding: 24, background: '#fff', textAlign: 'center' }}>
                <Route exact path='/' component={Push} />
                <Route path='/pop' component={Pop} />
                <Route path='/stock' component={Stock} />
                <Route path='/setting' component={Label} />
              </div>
            </Content>
            <Footer style={{ textAlign: 'center' }}>
        Ant Design ©2016 Created by Ant UED
            </Footer>
          </Layout>
        </Layout>
      </Router>
    )
  }
}

export default App
