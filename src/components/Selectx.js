import React, { Component } from 'react'
import axios from 'axios'
import { Select } from 'antd'
const Option = Select.Option

class Selectx extends Component {
  constructor (props) {
    super(props)
    this.state = {
      data: [],
      value: props.value || '',
      field: props.field || 'name'
    }
  }
  componentDidMount () {
    axios.get(this.props.url).then(
      res => this.setState({data: res.data})
    )
  }

  componentWillReceiveProps (nextProps) {
    // Should be a controlled component.
    const {field, value} = nextProps
    this.setState({field: field || 'name', value})
  }

  render () {
    const {data, value, field} = this.state
    return (
      <div>
        <Select value={value} onChange={this.props.onChange}>
          {
            data.map(item => (
              <Option key={item[field]}>{item.name}</Option>
            ))
          }
        </Select>
      </div>
    )
  }
}

export default Selectx
