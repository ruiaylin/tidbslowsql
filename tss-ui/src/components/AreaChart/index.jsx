import * as d3 from 'd3'

import React, { Component } from 'react'

import mockData from './mockData'

const duration = 1000

class AreaChart extends Component {
  constructor(props) {
    super(props)

    this.wrapper = React.createRef()
  }

  componentDidUpdate(prevProps) {
    if (this.props.data !== prevProps.data) {
      this.reRenderChart(this.props.data)
    }
  }

  componentDidMount() {
    this.initChart()
    this.reRenderChart(this.props.data && this.props.data.length > 0 ? this.props.data : mockData)
  }

  initChart = () => {
    const wrapper = this.wrapper.current

    const margin = { top: 15, right: 30, bottom: 30, left: 30 }
    const width = wrapper.offsetWidth - margin.left - margin.right
    const height = 250 - margin.top - margin.bottom
    this.height = height
    const color = this.props.color || '#00d1b2'

    const svg = d3
      .select(wrapper)
      .append('svg')
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .append('g')
      .attr('transform', `translate(${margin.left * 2}, ${margin.top})`)
    this.svg = svg

    const x = d3.scaleTime().range([0, width - margin.right * 2])
    this.x = x
    const xAxis = d3
      .axisBottom(x)
      .tickSize(0)
      .tickSizeOuter(0)
      .tickPadding(15)
    this.xAxis = xAxis
    svg
      .append('g')
      .attr('transform', `translate(0, ${height})`)
      .attr('class', 'area-chart-x')

    const y = d3.scaleLinear().range([height, 0])
    this.y = y
    const yAxis = d3
      .axisLeft(y)
      .tickSize(0)
      .tickSizeOuter(0)
      .tickPadding(15)
    this.yAxis = yAxis
    svg
      .append('g')
      .attr('transform', `translate(0, 0)`)
      .attr('class', 'area-chart-y')

    const focus = svg
      .append('g')
      .attr('class', 'focus')
      .style('display', 'none')

    focus.append('circle').attr('r', 5)

    const tooltip = d3
      .select('.tss-AreaChart')
      .append('div')
      .attr('id', 'area-chart-tooltip')
      .style('display', 'none')

    const tooltipDate = tooltip.append('div').attr('class', 'tooltip-date')
    const tooltipCount = tooltip.append('div').attr('class', 'tooltip-count')

    svg
      .append('path')
      .attr('class', 'area-chart-area')
      .attr('fill', color)
      .attr('fill-opacity', 0.6)
      .attr('stroke', 'none')
    svg
      .append('path')
      .attr('class', 'area-chart-line')
      .attr('fill', 'none')
      .attr('stroke', color)
      .attr('fill-opacity', 0.8)
      .attr('stroke-width', 3)

    svg
      .append('rect')
      .attr('fill', 'none')
      .attr('pointer-events', 'all')
      .attr('width', width - margin.right * 2)
      .attr('height', height)
      .on('mouseover', () => {
        focus.style('display', null)
        tooltip.style('display', null)
      })
      .on('mouseout', () => {
        focus.style('display', 'none')
        tooltip.style('display', 'none')
      })
      .on('mousemove', mousemove)

    const bisectDate = d3.bisector(d => d.date).left
    const dateFormatter = d3.timeFormat('%y/%m/%d-%H:%M')

    const that = this
    function mousemove() {
      if (!that.props.data) {
        return
      }
      const data = that.props.data && that.props.data.length > 0 ? that.props.data : mockData

      const x0 = x.invert(d3.mouse(this)[0])
      let i = bisectDate(data, x0)
      if (i === 0) {
        i = 1
      } else if (i === data.length) {
        i = data.length - 1
      }
      const d0 = data[i - 1]
      const d1 = data[i]
      const d = x0 - d0.date > d1.date - x0 ? d1 : d0

      focus.attr('transform', `translate(${x(d.date)}, ${y(d.count)})`)
      tooltip.attr('style', `top: ${y(d.count) - 64}px; left: ${x(d.date)}px;`)
      tooltipDate.text(dateFormatter(d.date))
      tooltipCount.text(`Count: ${d.count}`)
    }
  }

  reRenderChart = data => {
    this.x.domain(d3.extent(data, d => d.date))
    this.svg
      .select('.area-chart-x')
      .transition()
      .duration(duration)
      .call(this.xAxis)
    this.y.domain([0, ((d3.max(data, d => d.count) * 3) / 2).toFixed(0)])
    this.svg
      .select('.area-chart-y')
      .transition()
      .duration(duration)
      .call(this.yAxis)

    this.svg
      .select('.area-chart-area')
      .datum(data)
      .transition()
      .duration(duration)
      .attr(
        'd',
        d3
          .area()
          .x(d => this.x(d.date))
          .y0(this.height)
          .y1(d => this.y(d.count))
      )

    this.svg
      .select('.area-chart-line')
      .datum(data)
      .transition()
      .duration(duration)
      .attr(
        'd',
        d3
          .line()
          .x(d => this.x(d.date))
          .y(d => this.y(d.count))
      )
  }

  render() {
    return <div className="tss-AreaChart" ref={this.wrapper}></div>
  }
}

export default AreaChart
