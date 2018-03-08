import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm, InputNumber, Form, Modal } from 'antd'
import axios from 'axios'
import Selectx from './Selectx'

const FormItem = Form.Item

const BaseForm = Form.create()(
  (props) => {
    const { visible, onCancel, onOK, update, form, columns } = props
    const { getFieldDecorator } = form
    const formItemLayout = {
      labelCol: { span: 4 },
      wrapperCol: { span: 14 }
    }
    return (
      <Modal
        visible={visible}
        title={update ? '更新物品信息' : '创建物品'}
        okText={update ? '更新' : '创建'}
        onCancel={onCancel}
        onOk={onOK}
        cancelText='取消'
      >
        <Form layout='horizontal'>
          {
            columns.map(column => {
              const { name, type } = column
              let inner = null
              if (type === undefined || type === 'text') {
                inner = <Input />
              } else if (type === 'select') {
                inner = (
                  <Selectx
                    url={column.url}
                    initialValue={update ? update[name] : ''}
                  />
                )
              } else if (type === 'number') {
                inner = <InputNumber min={0} />
              }
              return (
                <FormItem label='名称' {...formItemLayout}>
                  {getFieldDecorator(name, {initialValue: update ? update[name] : ''})(
                    {inner}
                  )}
                </FormItem>
              )
            })
          }
        </Form>
      </Modal>
    )
  }
)

class Tablex extends Component {
  constructor (props) {
    super(props)

    this.state = {
      data: [],
      visible: false,
      update: null
    }
    const {columns, onDelete} = this.props
    const opertaion = {
      title: '操作',
      render: (text, record) => {
        return (
          <div>
            <Button
              type='primary'
              style={{ marginRight: '5px' }}
              onClick={() => this.setState({visible: true, update: record})}
            >
              编辑
            </Button>
            <Popconfirm
              placement='bottomLeft'
              title='确定删除吗？'
              onConfirm={() => onDelete(record.id)}
              okText='确定' cancelText='取消'
            >
              <Button>删除</Button>
            </Popconfirm>
          </div>
        )
      }
    }
    this.columns = [...columns, opertaion]
  }

  componentDidMount () {
    axios.get(`/${this.props.resource}`).then(
      res => {
        this.setState({data: res.data})
      }
    )
  }

  handleCancel = () => {
    this.setState({visible: false})
  }

  showModal = () => {
    this.setState({visible: true})
  }

  saveFormRef = (form) => {
    this.form = form
  }

  render () {
    const {update, visible} = this.state
    const {onCreate, onUpdate} = this.props
    return (
      <div>
        <Button
          className='editable-add-btn'
          onClick={() => this.setState({visible: true, update: null})}
        >
          添加
        </Button>
        <BaseForm
          ref={this.saveFormRef}
          visible={visible}
          onCancel={this.handleCancel}
          onOK={update ? onUpdate : onCreate}
          update={update}
          columns={this.props.columns}
        />
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
