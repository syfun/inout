import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm } from 'antd'
import axios from 'axios'

class Item extends Component {
  constructor (props) {
    super(props)
    this.state = { data: [] }
    this.columns = this.getColumns([
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
        'title': '入库',
        'dataIndex': 'push',
        'width': '5%'
      }, {
        'title': '出库',
        'dataIndex': 'pop',
        'width': '5%'
      }, {
        'title': '剩余',
        'dataIndex': 'now',
        'width': '5%'
      }
    ])
  }

  getColumns (options) {
    let columns = options.map(option => ({
      title: option.title,
      dataIndex: option.dataIndex,
      width: option.width,
      render: (text, record) => (
        <div>
          {record.editable
            ? <Input style={{ margin: '-5px 0' }} value={text} onChange={e => this.handleChange(e.target.value, record.id, option.dataIndex)} />
            : text
          }
        </div>
      )
    })
    )
    columns.push({
      title: 'operation',
      dataIndex: 'operation',
      render: (text, record) => {
        const { editable } = record
        return (
          <div className='editable-row-operations'>
            {
              editable
                ? <div>
                  <Button type='primary' onClick={() => this.save(record.id)}>保存</Button>
                  <Button onClick={() => this.cancel(record.id)}>取消</Button>
                </div>
                : <Button type='primary' onClick={() => this.edit(record.id)}>编辑</Button>
            }
            <Popconfirm
              placement='bottomLeft'
              title='确定删除吗？'
              onConfirm={() => this.delete(record.id)}
              okText='确定' cancelText='取消'
            >
              <Button>删除</Button>
            </Popconfirm>

          </div>
        )
      }
    })
    return columns
  }

  componentDidMount () {
    axios.get('/items').then(res => {
      this.setState({data: res.data})
      this.cacheData = res.data.map(item => ({ ...item }))
    })
  }

  handleChange (value, key, column) {
    const newData = [...this.state.data]
    const target = newData.filter(item => key === item.id)[0]
    if (target) {
      target[column] = value
      this.setState({ data: newData })
    }
  }
  edit (key) {
    const newData = [...this.state.data]
    const target = newData.filter(item => key === item.id)[0]
    if (target) {
      target.editable = true
      this.setState({ data: newData })
    }
  }
  save (key) {
    const newData = [...this.state.data]
    const target = newData.filter(item => key === item.id)[0]
    if (target) {
      axios.patch(`/items/${key}`, target).then(
        () => {
          delete target.editable
          this.setState({ data: newData })
          this.cacheData = newData.map(item => ({ ...item }))
        }
      )
    }
  }
  cancel (key) {
    const newData = [...this.state.data]
    const target = newData.filter(item => key === item.id)[0]
    if (target) {
      Object.assign(target, this.cacheData.filter(item => key === item.id)[0])
      delete target.editable
      this.setState({ data: newData })
    }
  }

  handleAdd () {
    const { data } = this.state
    const newCell = {
      id: data.length === 0 ? 1 : data[data.length - 1].id + 1
    }
    this.setState({data: [...data, newCell]})
  }

  delete (key) {
    axios.delete(`/items/${key}`).then(
      () => {
        let newData = [...this.state.data]
        newData = this.state.data.filter(item => key !== item.id)
        this.setState({ data: newData })
        this.cacheData = newData.map(item => ({ ...item }))
      }
    )
  }
  render () {
    return (
      <div>
        <Button className='editable-add-btn' onClick={this.handleAdd}>添加</Button>
        <Table
          bordered
          dataSource={this.state.data}
          columns={this.columns}
          rowKey='id'
        />
      </div>
    )
  }
}

export default Item
