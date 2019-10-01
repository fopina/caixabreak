import { h, Component } from 'preact';
import Button from 'preact-material-components/Button';
import 'preact-material-components/Button/style.css';
import style from './style';

export default class About extends Component {
	render() {
		return (
			<div class={`${style.profile} page`}>
				<h1>About</h1>
				<p>About stuff</p>
			</div>
		);
	}
}
