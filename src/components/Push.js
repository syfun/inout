import React from 'react'

import Tablex from './Tablex'

const columns = [
  {
    'title': '物品名称',
    'dataIndex': 'name',
    'width': '15%'
  }, {
    'title': '物品分类',
    'dataIndex': 'type',
    'width': '15%'
  }, {
    'title': '物品规格',
    'dataIndex': 'specification',
    'width': '15%'
  }, {
    'title': '单位',
    'dataIndex': 'unit',
    'width': '10%'
  }, {
    'title': '摘要',
    'dataIndex': 'abstract',
    'width': '10%'
  }, {
    'title': '入库数量',
    'dataIndex': 'number',
    'width': '10%',
    'type': 'number'
  }, {
    'title': '仓库',
    'dataIndex': 'warehouse',
    'width': '10%'
  }, {
    'title': '入库人',
    'dataIndex': 'user',
    'width': '10%'
  }, {
    'title': '备注',
    'dataIndex': 'remark',
    'width': '10%'
  }
]

function Push (props) {
  return (
    <Tablex columns={columns} resource='item' />
  )
}

export default Push
