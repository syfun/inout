import React, { Component } from 'react'
import { Table, Input, Button, Popconfirm, InputNumber, Form, Modal } from 'antd'
import Selectx from './Selectx'

const FormItem = Form.Item

const BaseForm = Form.create()(
  (props) => {
    const { visible, onCancel, onOK, update, form } = props
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
          <FormItem label='名称' {...formItemLayout}>
            {getFieldDecorator('name', {initialValue: update ? update.name : ''})(
              <Input />
            )}
          </FormItem>
          <FormItem label='物品分类' {...formItemLayout}>
            {getFieldDecorator('type')(
              <Selectx
                url='/labels?type=itemType'
                initialValue={update ? update.type : ''}
              />
            )}
          </FormItem>
          <FormItem label='物品规格' {...formItemLayout}>
            {getFieldDecorator('specification', {initialValue: update ? update.specification : ''})(
              <Input />
            )}
          </FormItem>
          <FormItem label='单位' {...formItemLayout}>
            {getFieldDecorator('unit', {initialValue: update ? update.unit : ''})(
              <Input />
            )}
          </FormItem>
          <FormItem label='入库' {...formItemLayout}>
            {getFieldDecorator('push', {initialValue: update ? update.push : 0})(
              <InputNumber min={0} />
            )}
          </FormItem>
          <FormItem label='出库' {...formItemLayout}>
            {getFieldDecorator('pop', {initialValue: update ? update.pop : 0})(
              <InputNumber min={0} />
            )}
          </FormItem>
          <FormItem label='剩余' {...formItemLayout}>
            {getFieldDecorator('now', {initialValue: update ? update.now : 0})(
              <InputNumber min={0} />
            )}
          </FormItem>
        </Form>
      </Modal>
    )
  }
)

class Item extends Component {
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

export default Item
