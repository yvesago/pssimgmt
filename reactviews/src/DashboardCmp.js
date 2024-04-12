import React, { useState, useEffect } from 'react';
import { TreeView, TreeItem } from '@mui/x-tree-view';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

import { useStore, useDataProvider, Title } from 'react-admin';


import { MyConfig } from './MyConfig';

const classes = {
    height: 110,
    flexGrow: 1,
    maxWidth: 400,
};


const renderTree = (nodes) => (
    <TreeItem key={nodes.id} nodeId={nodes.id.toString()} label={nodes.name}>
        {Array.isArray(nodes.children) ? nodes.children.map((node) => renderTree(node)) : null}
    </TreeItem>
);

let leaf = {};
const DashboardCmp = () => {
    const dataProvider = useDataProvider();
    const [state, setState] = useState({data: null, parents: ['0'], dom: 0});
    const [dom, setDom] = useStore('dom', {select: {id:0 }, evaluation: {}});

    //export default function RecursiveTreeView() {
    //state = {data: null, parents: ["0"], dom: 0};
    //dom = {select : {id: "0"} };
    //parents = [0];
  
    const findleaf = (n, data) => {
        if (n) {
            if ( Array.isArray(n.children) === true) { 
                const p = state.parents;
                p.push(n.id.toString());
                setState({data: data, parents : p, dom: state.dom });
                n.children.map((node) => findleaf(node ,data)); 
            }
            else {
                leaf[n.id] = true;
            }
        }
    };


    useEffect(() => {
        if (dom === undefined) { return; }
        //console.log('====CMP======');
        dataProvider.getOne('themestree', {id: dom.select.id})
            .then(({ data }) => {
                //console.log('--dataProvider--');
                //console.log(data);
                setDom({
                    select : (dom) ? dom.select: {},
                    mask : (dom) ? dom.mask: {},
                    evaluation: data.evaluation,
                });

                setState({data : data});
                findleaf(data,data);
            })
            .catch(error => {
                //setError(error);
                console.log('-error-');
                console.log(error);
                window.location.href = MyConfig.BASE_PATH + '#/login';
                return;
                //setLoading(false);
            });
    },[dom.select.id]);


    const selectNode = (e,nodeid) => { 
        //console.log("--" + nodeid);
        //console.log(leaf[nodeid]);
        if ( leaf[nodeid] === true ) {
            var redirect = '/themes/'+nodeid+'/show';
            window.location.href = MyConfig.BASE_PATH + '#' + redirect;
        }
    };

    if (dom === undefined || !state.data || leaf === null) {
        return <div>...</div>;
    } else {
        return (
            <>
                <Title title={MyConfig.NAME} />
                {dom.select.id !== undefined && dom.select.id !== '0' ?
                    <h2>Périmètre : {dom.select.name}</h2>
                    : null }
                <TreeView
                    className={classes.root}
                    defaultCollapseIcon={<ExpandMoreIcon />}
                    defaultExpanded={state.parents}
                    defaultExpandIcon={<ChevronRightIcon />}
                    onNodeSelect={selectNode}
                >
                    {renderTree(state.data)}
                </TreeView>
            </>
        );
    }
};

export default DashboardCmp;
