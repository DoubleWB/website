import React from 'react';
import HomeComponent from './components/HomeComponent'
import CascaleComponent from './components/CascaleComponent'
import AboutComponent from './components/AboutComponent'
import { BrowserRouter as Router, Link, Route } from 'react-router-dom'

export default class App extends React.Component {
    render() {
        return (
            <Router>
                <div className="container">
                    <h1>Will B. </h1>
                    <h2>Robotics // AI // Devops </h2>
                    <Link to="/">Home</Link>
                    &nbsp;
                    <Link to="/cascale">CaScale (Scale Project)</Link>
                    &nbsp;
                    <Link to="/about">About</Link>
                    &nbsp;
                    <p></p>
                    <Route
                        exact path="/"
                        render={() => <HomeComponent />} />
                    <Route
                        exact path="/cascale"
                        render={() => <CascaleComponent />} />
                    <Route
                        exact path="/about"
                        render={() => <AboutComponent />} />
                </div>
            </Router>
        )
    }
}
