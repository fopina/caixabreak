import { h, Component } from 'preact';
import Button from 'preact-material-components/Button';
import TextField from 'preact-material-components/TextField';
import linkState from 'linkstate';
import LinearProgress from 'preact-material-components/LinearProgress';
import Dialog from 'preact-material-components/Dialog';
import Checkbox from 'preact-material-components/Checkbox'
import Formfield from 'preact-material-components/FormField';
import 'preact-material-components/Checkbox/style.css';
import 'preact-material-components/Button/style.css';
import 'preact-material-components/TextField/style.css';
import 'preact-material-components/LinearProgress/style.css';
import 'preact-material-components/Dialog/style.css';
import style from './style';
import {setToken, saveLogin, removeLogin, refreshWithLogin, updateData} from '../../utils/BreakService'
import { route } from 'preact-router';

export default class Login extends Component {
	state = {
		loading: false,
		username: '',
		password: '',
		rememberMe: true,
	};

	login = () => {
		console.log(this.state.rememberMe)
		this.setState({ loading: true })
		refreshWithLogin(this.state.username, this.state.password).then((data) => {
			this.setState({ loading: false })
			if (this.state.rememberMe) {
				saveLogin(this.state.username, this.state.password)
			}
			setToken(data.Token)
			updateData(data)
			route("/", true)
		}).catch((error) => {
			this.setState({ loading: false })
			if ((error.response) && (error.response.status == 401)) {
				removeLogin()
				this.scrollingDlg.MDComponent.show()
			} else {
				alert('unexpected error')
			}
		});
	};

	keypress = e => {
		if(e.key === 'Enter'){
            this.login()
        }
	};
	
	render(props, state) {
		if (state.loading) {
			return (
				<div class={`${style.profile} page`}>
				<h4>Logging in...</h4>
					<LinearProgress indeterminate />
					<Dialog ref={scrollingDlg=>{this.scrollingDlg=scrollingDlg;}}>
					</Dialog>
				</div>
			)
		}
		return (
			<div class={`${style.profile} page`}>
				<h1>Login</h1>
				<Formfield>
					<TextField label="Username" onKeyPress={this.keypress} value={state.username} onInput={linkState(this, 'username')}/>
				</Formfield>
				<Formfield>
					<TextField label="Password" type="password" onKeyPress={this.keypress} value={state.password} onInput={linkState(this, 'password')}/>
				</Formfield>
				<Formfield>
					<Checkbox id="basic-checkbox" class={style.center_checkbox} checked={state.rememberMe} onChange={linkState(this, 'rememberMe')}/>
					<label for="basic-checkbox">Remember me</label>
				</Formfield>
				<p>
					<Button raised ripple onClick={() => this.login()}>Log in</Button>
				</p>
				<Dialog ref={scrollingDlg=>{this.scrollingDlg=scrollingDlg;}}>
					<Dialog.Header>Error</Dialog.Header>
					<Dialog.Body scrollable={true}>
						Invalid login
					</Dialog.Body>
					<Dialog.Footer>
						<Dialog.FooterButton accept={true}>OK</Dialog.FooterButton>
					</Dialog.Footer>
				</Dialog>
			</div>
		);
	}
}
