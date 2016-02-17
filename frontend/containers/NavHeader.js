import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'
import { push } from 'react-router-redux'
import { loadAuthUser, logout } from '../actions/auth'
import isString from 'lodash/isString'
import bindAll from 'lodash/bindAll'
import Radium from 'radium'
import styles from '../styles'

class NavHeader extends Component {
  componentWillMount() {
    bindAll(this, 'handleLogout')
  }

  handleLogout(ev) {
    ev.preventDefault()
    this.props.logout()
    this.props.push('/')
  }

  render() {
    const { activeTab, isLoggedIn, authUser } = this.props
    const homeActiveClass = activeTab === 'home' ? 'active' : ''
    const createModelActiveClass = activeTab === 'create-model' ? 'active' : ''
    const profileActiveClass = activeTab === 'profile' ? 'active' : ''
    const loginActiveClass = activeTab === 'login' ? 'active' : ''
    const registerActiveClass = activeTab === 'register' ? 'active' : ''
    return (
      <div className="header clearfix" style={styles.header}>
        <nav>
          <ul className="nav nav-pills pull-right">
            <li className={homeActiveClass} htmlRole="presentation">
              <Link to="/">Home</Link>
            </li>
            {isLoggedIn ?
              <li className={createModelActiveClass} htmlRole="presentation">
                <Link to="/create-model">Create Model</Link>
              </li> : null}
            {authUser ?
              <li className={profileActiveClass} htmlRole="presentation">
                <Link to={'/' + authUser.username}>{authUser.username}</Link>
              </li> : null}
            <li className={loginActiveClass} htmlRole="presentation">
              {isLoggedIn ? null :
                <Link to="/login">Login</Link>}
            </li>
            <li className={registerActiveClass} htmlRole="presentation">
              {isLoggedIn ?
                <a href="#" onClick={this.handleLogout}>Logout</a> :
                <Link to="/register">Sign Up</Link>}
            </li>
          </ul>
        </nav>
        <h3 className="text-muted" style={styles.masthead}>Gradientzoo</h3>
      </div>
    )
  }
}

NavHeader.PropTypes = {
  activeTab: PropTypes.string,
  isLoggedIn: PropTypes.bool.isRequired,
  authUser: PropTypes.object,
  logout: PropTypes.func.isRequired,
  push: PropTypes.func.isRequired
}

function mapStateToProps(state, props) {
  return {
    isLoggedIn: isString(state.authUserId) && state.authUserId.length > 0,
    authUser: state.authUserId ? state.entities.users[state.authUserId] : null
  }
}

export default Radium(connect(mapStateToProps, {
  logout,
  push
})(NavHeader))
