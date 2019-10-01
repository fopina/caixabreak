import { Component } from 'preact';
import List from 'preact-material-components/List';
import Fab from 'preact-material-components/Fab';
import LinearProgress from 'preact-material-components/LinearProgress';
import ListItemMetaText from '../../components/material-listmetatext'
import 'preact-material-components/List/style.css';
import 'preact-material-components/Fab/style.css';
import 'preact-material-components/LinearProgress/style.css';
import style from './style';
import { getData, updateData, getToken, setToken, refreshWithToken } from '../../utils/BreakService'
import { route } from 'preact-router';

export default class Home extends Component {
	state = {
		loading: false
	}

	refresh = () => {
		this.setState({ loading: true })
		refreshWithToken(getToken()).then((data) => {
			updateData(data)
			this.setState({ loading: false })
		}).catch((error) => {
			this.setState({ loading: false })
			console.log(error)
			if (error.response && error.response.data) {
				if(error.response.data.Error === "not logged in") {
					setToken(null)
					route("/login/", true)
					return
				}
				else if(error.response.data.Error === "invalid token") {
					setToken(null)
					route("/login/", true)
					return
				}
			}
			alert('unexpected error')
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
					<List two-line={true}>
					{ getData().History.map(item => 
							<List.Item>
								<List.TextContainer>
									<List.PrimaryText>{item.Description}</List.PrimaryText>	
									<List.SecondaryText>{item.Date}</List.SecondaryText>	
								</List.TextContainer>
								<ListItemMetaText class={
									item.CreditAmount > 0 
									? `${style.credit}` 
									: `${style.debit}`
									}>
								{
									item.CreditAmount > 0
									? item.CreditAmount 
									: item.DebitAmount
								}
								</ListItemMetaText>
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
