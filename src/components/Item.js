import React from 'react'

import Tablex from './Tablex'

const columns = [
  {
    title: '物品名称',
    dataIndex: 'name',
    width: '15%'
  }, {
    title: '物品分类',
    dataIndex: 'type',
    width: '15%',
    type: 'select',
    url: '/labels?type=itemType'
  }, {
    title: '物品规格',
    dataIndex: 'specification',
    width: '15%'
  }, {
    title: '单位',
    dataIndex: 'unit',
    width: '10%'
  }, {
    title: '入库',
    dataIndex: 'push',
    width: '10%',
    type: 'number'
  }, {
    title: '出库',
    dataIndex: 'pop',
    width: '10%',
    type: 'number'
  }, {
    title: '剩余',
    dataIndex: 'now',
    width: '10%',
    type: 'number'
  }
]

const Item = (props) => (
  <Tablex columns={columns} resource='item' />
)

export default Item
