import * as actions from './'
import * as realtimeAPI from '../api/realtime'
import * as reportAPI from '../api/report'
import * as sqlAPI from '../api/sql'

function getRealtimeTrending(params) {
  return dispatch => {
    realtimeAPI.getRealtimeTrending(params).then(resp => {
      const data = resp.data.results

      dispatch(actions.realtimeTrendingGet(data))
    })
  }
}

function getRealtimeQuery(params) {
  return dispatch => {
    realtimeAPI.getRealtimeQuery(params).then(resp => {
      const data = resp.data.results

      dispatch(actions.realtimeQueryGet(data))
    })
  }
}

function getReport(params) {
  return dispatch => {
    reportAPI.getReport(params).then(resp => {
      const data = resp.data.results

      dispatch(actions.reportGet(data))
    })
  }
}

function getSqlInfo(cluster, queryId) {
  return dispatch => {
    sqlAPI.getSqlInfo(cluster, queryId).then(resp => {
      const data = resp.data.results

      dispatch(actions.sqlInfoSet(data))
    })
  }
}

function getSingleSqlTrending(cluster, sqlId) {
  return dispatch => {
    sqlAPI.getSingleSqlTrending(cluster, sqlId).then(resp => {
      const data = resp.data.results

      dispatch(actions.singleSqlTrendingSet(data))
    })
  }
}

export { getRealtimeTrending, getRealtimeQuery, getReport, getSqlInfo, getSingleSqlTrending }
