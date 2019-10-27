import { getRealtimeQuery, getRealtimeTrending, getReport } from '../../actions/apis'
import { realtimeClusterSet, realtimeInstanceSet, reportClusterSet, reportInstanceSet } from '../../actions'

import DateTimePicker from '../DateTimePicker'
import React, { useState } from 'react'
import { connect } from 'react-redux'

const TssInputGroup = props => {
  const [params, setParams] = useState({
    cluster: 'ugc',
    instance: '',
    start: '',
    end: ''
  })

  const handleInputChange = type => e => {
    setParams({ ...params, [type]: e.target.value })
  }

  const handleInputOnKeyDown = e => {
    if (e.key === 'Enter') {
      handleDifferentGet()
    }
  }

  const handleDifferentGet = () => {
    switch (props.type) {
      case 'realtime':
        props.handleSetRealtimeCluster(params.cluster)
        props.handleSetRealtimeInstance(params.instance)
        props.handleGetRealtimeTrending(params)
        props.handleGetRealtimeQuery(params)
        break
      case 'report':
        props.handleSetReportCluster(params.cluster)
        props.handleSetReportInstance(params.instance)
        props.handleGetReport(params)
        break
      default:
        return
    }
  }

  const handlePickCallback = (start, end) => {
    setParams({ ...params, start, end })
  }

  return (
    <article className="panel tss-TssInputGroup">
      <div className="panel-block">
        <div className="field is-horizontal">
          <div className="field-label">
            <label className="label">Cluster</label>
          </div>
          <div className="field">
            <p className="control">
              <input
                className="input"
                type="text"
                placeholder="Type a cluster name"
                value={params.cluster}
                onChange={handleInputChange('cluster')}
                onKeyDown={handleInputOnKeyDown}
              />
            </p>
          </div>
          <div className="field-label">
            <label className="label">Instance</label>
          </div>
          <div className="field">
            <p className="control">
              <input
                className="input"
                type="text"
                placeholder="Type a instance name"
                onChange={handleInputChange('instance')}
                onKeyDown={handleInputOnKeyDown}
              />
            </p>
          </div>
          {props.type === 'report' && (
            <div className="field-label">
              <label className="label">Date</label>
            </div>
          )}
          {props.type === 'report' && <DateTimePicker handlePickCallback={handlePickCallback} />}
          <div className="field field-submit-button">
            <div className="control">
              <button className="button is-primary" onClick={handleDifferentGet}>
                Submit
              </button>
            </div>
          </div>
        </div>
      </div>
    </article>
  )
}

const mapDispatchToProps = dispatch => ({
  handleGetRealtimeTrending: params => dispatch(getRealtimeTrending(params)),
  handleGetRealtimeQuery: params => dispatch(getRealtimeQuery(params)),
  handleSetRealtimeCluster: cluster => dispatch(realtimeClusterSet(cluster)),
  handleSetRealtimeInstance: instance => dispatch(realtimeInstanceSet(instance)),
  handleGetReport: params => dispatch(getReport(params)),
  handleSetReportCluster: cluster => dispatch(reportClusterSet(cluster)),
  handleSetReportInstance: instance => dispatch(reportInstanceSet(instance))
})

export default connect(
  null,
  mapDispatchToProps
)(TssInputGroup)
