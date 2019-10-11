import { Component } from 'preact';
import { route } from 'preact-router';
import TopAppBar from 'preact-material-components/TopAppBar';
import Drawer from 'preact-material-components/Drawer';
import List from 'preact-material-components/List';
import Dialog from 'preact-material-components/Dialog';
import Button from 'preact-material-components/Button';
import 'preact-material-components/Dialog/style.css';
import 'preact-material-components/Drawer/style.css';
import 'preact-material-components/List/style.css';
import 'preact-material-components/TopAppBar/style.css';
import 'preact-material-components/Button/style.css';
import { isLoggedIn, getLoggedUser, logout } from '../../utils/BreakService';

export default class Header extends Component {
	closeDrawer() {
		this.drawer.MDComponent.open = false;
	}

	openDrawer = () => (this.drawer.MDComponent.open = true);

	openSettings = () => this.dialog.MDComponent.show();

	drawerRef = drawer => (this.drawer = drawer);
	dialogRef = dialog => (this.dialog = dialog);

	linkTo = path => () => {
		route(path);
		this.closeDrawer();
	};

	goHome = this.linkTo('/');
	goToAbout = this.linkTo('/about');
	version = process.env.PREACT_APP_VERSION || 'dev'

	loggedUser = () => {
		return getLoggedUser()
	}

	logout = () => {
		logout()
		this.dialog.MDComponent.close()
		route('/login/')
	}

	reload = () => {
		window.location.reload();
	}

	render(props) {
		return (
			<div>
				<TopAppBar fixed={true} className="topappbar">
					<TopAppBar.Row>
						<TopAppBar.Section align-start>
							<TopAppBar.Icon menu onClick={this.openDrawer}>
								menu
							</TopAppBar.Icon>
							<TopAppBar.Title>Break</TopAppBar.Title>
						</TopAppBar.Section>
						<TopAppBar.Section align-end shrink-to-fit onClick={this.openSettings}>
							<TopAppBar.Icon>settings</TopAppBar.Icon>
						</TopAppBar.Section>
					</TopAppBar.Row>
				</TopAppBar>
				<Drawer modal ref={this.drawerRef}>
					<Drawer.DrawerContent>
						<Drawer.DrawerItem tabindex={0} selected={props.selectedRoute === '/'} onClick={this.goHome}>
							<List.ItemGraphic>home</List.ItemGraphic>
							Home
						</Drawer.DrawerItem>
						<Drawer.DrawerItem tabindex={0} selected={props.selectedRoute === '/about'} onClick={this.goToAbout}>
							<List.ItemGraphic>info</List.ItemGraphic>
							About
						</Drawer.DrawerItem>
					</Drawer.DrawerContent>
				</Drawer>
				<Dialog ref={this.dialogRef}>
					<Dialog.Header>Settings</Dialog.Header>
					<Dialog.Body>
					{
						isLoggedIn() &&
						<div>
							Account: {this.loggedUser()}
						</div>
					}
						<div>
							Version: {this.version}
						</div>
					</Dialog.Body>
					<Dialog.Body>
						<div>
						{
							isLoggedIn() &&
							<Button raised={true} onClick={this.logout} class="left-button">Logout</Button>
						}
							<Button raised={true} onClick={this.reload}>Reload app</Button>
						</div>
					</Dialog.Body>
					<Dialog.Footer>
						<Dialog.FooterButton accept>OK</Dialog.FooterButton>
					</Dialog.Footer>
				</Dialog>
			</div>
		);
	}
}
