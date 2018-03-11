import React, { Component } from 'react'
import { Menu, Icon } from 'antd'

import { NavLink } from 'react-router-dom'

class LeftMenu extends Component {
  render () {
    return (
      <Menu
        defaultSelectedKeys={['1']}
        defaultOpenKeys={['1']}
        mode='horizontal'
        theme='dark'
        style={{ lineHeight: '64px' }}
      >
        <Menu.Item key='1'>
          <NavLink to='/push'>
            <Icon type='pie-chart' />
            <span>入库</span>
          </NavLink>

        </Menu.Item>
        <Menu.Item key='2'>
          <NavLink to='/pop'>
            <Icon type='desktop' />
            <span>出库</span>
          </NavLink>
        </Menu.Item>
        <Menu.Item key='3'>
          <NavLink to='/stock'>
            <Icon type='inbox' />
            <span>库存统计</span>
          </NavLink>
        </Menu.Item>
        <Menu.Item key='4'>
          <NavLink to='/item'>
            <Icon type='inbox' />
            <span>物品</span>
          </NavLink>
        </Menu.Item>
        <Menu.Item key='5'>
          <NavLink to='/setting'>
            <Icon type='inbox' />
            <span>设置</span>
          </NavLink>
        </Menu.Item>
      </Menu>
    )
  }
}

export default LeftMenu
