import MaterialComponent from 'preact-material-components/Base/MaterialComponent';

export default class ListItemMetaText extends MaterialComponent {
    componentName = 'list-item__meta'
    mdcProps = []

    materialDom(props) {
        return (
          <span {...props} ref={this.setControlRef} role="presentation">
            {props.children}
          </span>
        );
    }
}
