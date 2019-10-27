import React, { useEffect, useRef, useState } from 'react'
import { getRealtimeQuery, getRealtimeTrending, getSqlInfo, getSingleSqlTrending } from '../../actions/apis'

import AreaChart from '../AreaChart'
import Table from '../Table'
import TableHeader from '../Table/TableHeader'
import TableItem from '../Table/TableItem'
import TssInputGroup from '../TssInputGroup'
import { connect } from 'react-redux'

const Realtime = props => {
  const realtimeId = useRef(null)
  const [realtime, setRealtime] = useState(false)

  useEffect(
    () => () => {
      clearInterval(realtimeId.current)
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    []
  )

  const handleRealtime = () => {
    if (realtime) {
      setRealtime(false)
      clearInterval(realtimeId.current)
    } else {
      setRealtime(true)
      realtimeId.current = setInterval(() => {
        props.handleGetRealtimeTrending({
          cluster: props.realtimeCluster,
          instance: props.realtimeInstance,
          start: '',
          end: ''
        })
        props.handleGetRealtimeQuery({
          cluster: props.realtimeCluster,
          instance: props.realtimeInstance,
          start: '',
          end: ''
        })
      }, 5000)
    }
  }

  return (
    <div className="tss-Realtime">
      <TssInputGroup type="realtime" />
      <article className="panel panel-trending">
        <div className="panel-heading tss-panel-heading">
          <div className="title is-4">
            Realtime Trending {props.realtimeCluster ? `(Cluster: ${props.realtimeCluster})` : ''}{' '}
            {props.realtimeInstance ? `(Instance: ${props.realtimeInstance})` : ''}
          </div>
          {props.realtimeCluster && (
            <div onClick={handleRealtime}>
              <i
                style={{
                  marginBottom: '1rem',
                  cursor: 'pointer'
                }}
                className={`fas fa-circle-notch has-text-primary${realtime ? ' fa-spin' : ''}`}
              />
            </div>
          )}
        </div>
        <AreaChart data={props.trendingData} />
      </article>
      <article className="panel panel-query">
        <div className="panel-heading tss-panel-heading">
          <div className="title is-4">Realtime Query</div>
        </div>
        <Table
          header={TableHeader}
          item={TableItem}
          data={props.queryData}
          handleGetSqlInfo={props.handleGetSqlInfo}
          handleGetSingleSqlTrending={props.handleGetSingleSqlTrending}
        />
      </article>
    </div>
  )
}

const mapStateToProps = state => ({
  trendingData: state.realtimeTrending.map(d => ({
    date: new Date(d.day_time),
    count: d.counter
  })),
  queryData: state.realtimeQuery,
  realtimeCluster: state.realtimeCluster,
  realtimeInstance: state.realtimeInstance
})

const mapDispatchToProps = dispatch => ({
  handleGetRealtimeTrending: params => dispatch(getRealtimeTrending(params)),
  handleGetRealtimeQuery: params => dispatch(getRealtimeQuery(params)),
  handleGetSqlInfo: (cluster, queryId) => dispatch(getSqlInfo(cluster, queryId)),
  handleGetSingleSqlTrending: (cluster, sqlId) => dispatch(getSingleSqlTrending(cluster, sqlId))
})

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Realtime)
