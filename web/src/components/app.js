import { Component } from 'preact';
import { Router, route } from 'preact-router';

import Header from './header';
import Home from '../routes/home';
import About from '../routes/about';
import Login from '../routes/login';
import NotFound from '../routes/404';
import { isLoggedIn } from '../utils/BreakService';

export default class App extends Component {
	handleRoute = e => {
		const isAuthed = isLoggedIn();
		if ((e.current.attributes?e.current.attributes.auth:false) && !isAuthed) {
			route("/login/", true)
		}
		this.setState({
			currentUrl: e.url
		});
	};

	render() {
		return (
			<div id="app">
				<Header selectedRoute={this.state.currentUrl} />
				<Router onChange={this.handleRoute}>
					<Home path="/" auth={true} />
					<About path="/about/" />
					<Login path="/login/" />
					<NotFound default />
				</Router>
			</div>
		);
	}
}
