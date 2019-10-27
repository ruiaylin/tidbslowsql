import React from 'react'
import Table from '../Table'
import TableHeader from '../Table/TableHeader'
import TableItem from '../Table/TableItem'
import TssInputGroup from '../TssInputGroup'
import { connect } from 'react-redux'
import { getSqlInfo, getSingleSqlTrending } from '../../actions/apis'

const Report = props => {
  return (
    <div className="tss-Report">
      <TssInputGroup type="report" />
      <article className="panel panel-report">
        <div className="panel-heading tss-panel-heading">
          <div className="title is-4">
            Statistical Report {props.reportCluster ? `(Cluster: ${props.reportCluster})` : ''}{' '}
            {props.reportInstance ? `(Instance: ${props.reportInstance})` : ''}
          </div>
        </div>
        <Table
          header={TableHeader}
          item={TableItem}
          data={props.reportData}
          handleGetSqlInfo={props.handleGetSqlInfo}
          handleGetSingleSqlTrending={props.handleGetSingleSqlTrending}
        />
      </article>
    </div>
  )
}

const mapStateToProps = state => ({
  reportData: state.report,
  reportCluster: state.reportCluster,
  reportInstance: state.reportInstance
})

const mapDispatchToProps = dispatch => ({
  handleGetSqlInfo: (cluster, queryId) => dispatch(getSqlInfo(cluster, queryId)),
  handleGetSingleSqlTrending: (cluster, sqlId) => dispatch(getSingleSqlTrending(cluster, sqlId))
})

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Report)
