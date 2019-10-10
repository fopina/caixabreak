import { Component } from 'preact';
import List from 'preact-material-components/List';
import Fab from 'preact-material-components/Fab';
import LinearProgress from 'preact-material-components/LinearProgress';
import Card from 'preact-material-components/Card';
import Icon from 'preact-material-components/Icon';
import 'preact-material-components/Card/style.css';
import 'preact-material-components/Icon/style.css';
import 'preact-material-components/List/style.css';
import 'preact-material-components/Fab/style.css';
import 'preact-material-components/LinearProgress/style.css';
import style from './style';
import { getData, updateData, setToken, refreshWithSavedLogin, refreshWithToken } from '../../utils/BreakService'
import { route } from 'preact-router';

export default class Home extends Component {
	state = {
		loading: false
	}

	refresh = () => {
		this.setState({ loading: true })
		refreshWithToken().then((data) => {
			updateData(data)
			this.setState({ loading: false })
			return
		}).catch((error) => {
			console.log(error)
			if ((error.response) && (error.response.status == 401)) {
				refreshWithSavedLogin().then((data) => {
					setToken(data.Token)
					updateData(data)
					this.setState({ loading: false })
					return
				}).catch((error) => {
					console.log(error)
					if ((error.response) && (error.response.status == 401)) {
						setToken(null)
						route("/login/", true)
						return
					} else {
						this.setState({ loading: false })
						alert('unexpected error')
					}
				})
			} else {
				this.setState({ loading: false })
				alert('unexpected error')
			}
		});
	}

	render(props, state) {
		if (state.loading) {
			return (
				<div class={`${style.home} page`}>
				<h4>Refreshing...</h4>
					<LinearProgress indeterminate />
				</div>
			)
		} else {
			return (<div class={`${style.home} page`}>
				<h1></h1>
				<div>
						<Card>
							<div class={style.cardHeader}>
							<div class="mdc-typography--title">Balance</div>
								<h2 class="mdc-typography--caption">
									<Icon>euro_symbol</Icon>
									{getData().Balance}
								</h2>
							</div>
						</Card>
					</div>
				<div>
					<List two-line={true}>
					{ getData().History.map(item => 
							<List.Item>
								<List.TextContainer>
									<List.PrimaryText>{item.Description}</List.PrimaryText>	
									<List.SecondaryText>{item.Date}</List.SecondaryText>	
								</List.TextContainer>
								<List.ItemMetaText class={
									item.CreditAmount > 0 
									? `${style.credit}` 
									: `${style.debit}`
									}>
								{
									item.CreditAmount > 0
									? item.CreditAmount 
									: item.DebitAmount
								}
								</List.ItemMetaText>
							</List.Item>
					)}
					</List>
				</div>
				<Fab ripple class={style.fab} onClick={this.refresh}><Fab.Icon>refresh</Fab.Icon></Fab>
			</div>
		)
		}
	}
}
