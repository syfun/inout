import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm, InputNumber } from 'antd'
import axios from 'axios'
import Selectx from './Selectx'

const EditableCell = ({ editable, value, onChange, option }) => {
  const type = option.type === undefined ? 'text' : option.type
  let editor = null
  if (type === 'number') {
    editor = <InputNumber min={0} defaultValue={value} onChange={onChange} />
  } else if (type === 'select') {
    editor = <Selectx url={option.url} defaultValue={value} onChange={onChange} />
  } else {
    editor = <Input style={{ margin: '-5px 0' }} type={type} value={value} onChange={e => onChange(e.target.value)} />
  }
  return (
    <div>
      {editable ? editor : value}
    </div>
  )
}

class Tablex extends Component {
  constructor (props) {
    super(props)
    this.columns = this.getColumns(props.columns)
    this.state = {
      data: [],
      resource: props.resource,
      query: props.query,
      extraData: props.extraData
    }
  }

  getColumns (options) {
    let columns = options.map(option => ({
      title: option.title,
      dataIndex: option.dataIndex,
      width: option.width,
      render: (text, record) => this.renderColumn(text, record, option)
    }))
    columns.push({
      title: '操作',
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
    const { resource, query } = this.state
    let url = `/${resource}s`
    if (query !== undefined) {
      url += `?${query}`
    }
    axios.get(url).then(res => {
      this.setState({data: res.data})
      this.cacheData = res.data.map(item => ({ ...item }))
    })
  }

  renderColumn = (text, record, option) => (
    <EditableCell
      editable={record.editable}
      value={text}
      option={option}
      onChange={value => this.handleChange(value, record.id, option.dataIndex)}
    />
  )

  handleChange (value, key, column) {
    console.log(value)
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
      if (target.add !== undefined && target.add) {
        axios.post(`/${this.state.resource}s`, target).then(
          res => {
            Object.assign(target, res.data)
            delete target.editable
            this.setState({data: newData})
            this.cacheData = newData.map(item => ({ ...item }))
          }
        )
      } else {
        axios.patch(`/${this.state.resource}s/${key}`, target).then(
          () => {
            delete target.editable
            this.setState({ data: newData })
            this.cacheData = newData.map(item => ({ ...item }))
          }
        )
      }
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

  handleAdd = () => {
    const { data, extraData } = this.state
    let newCell = {
      id: data.length === 0 ? 1 : data[data.length - 1].id + 1,
      editable: true,
      add: true
    }
    if (extraData !== undefined) {
      Object.assign(newCell, extraData)
    }
    this.setState({data: [...data, newCell]})
  }

  delete (key) {
    axios.delete(`/${this.state.resource}s/${key}`).then(
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

export default Tablex
