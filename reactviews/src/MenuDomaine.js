import React, { useEffect, useState } from 'react';
import { TreeView, TreeItem } from '@mui/x-tree-view';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { useStore, useDataProvider } from 'react-admin';

import { MyConfig } from './MyConfig';


const classes = { height: 110, flexGrow: 1, maxWidth: 400 };

const renderTree = (nodes) => (
    <TreeItem key={nodes.id} nodeId={nodes.id.toString()} label={nodes.name}>
        {Array.isArray(nodes.children) ? nodes.children.map((node) => renderTree(node)) : null}
    </TreeItem>
);


const MenuDomaine = () => {
    const [state, setState] = useState({data: null, parents: ['0'], lobj: {}});
    const [dom, setDom] = useStore('dom', {select: {id:0 }, evaluation: {}});
    const dataProvider = useDataProvider();

    let leaf = {};

    const findparents = (n, data) => {
        if (n) {
            const lobj = state.lobj;
            if ( lobj[n] !== 0 && lobj[n] !== undefined ) {
                const p = state.parents;
                p.push(lobj[n].id.toString());
                setState({data: data, lobj: state.lobj, parents : p});
                findparents(lobj[n].parent, data);
            } 
        }
    };

    const findleaf = (n) => {
        if (n) {
            if ( Array.isArray(n.children) === true) {
                const nsort = n.children.sort((a, b) => a.name.localeCompare(b.name));
                nsort.map((node) => findleaf(node));
            }
            else {
                leaf[n.id] = true;
            }
            const lobj = state.lobj;
            lobj[n.id]=n;
            setState({data: n, lobj: lobj, parents: state.parents});
        }

    };

    useEffect(() => {
        //console.log('=====Menu=====');
        dataProvider.getOne('', {id:'domainestree'})
            .then(({ data }) => {
                //console.log('--dataProvider--');
                //console.log(data);
                setState({data : data});
                findleaf(data);
                if (dom !== undefined && dom.select.id) {
                    findparents(dom.select.id, data);
                }
            })
            .catch(error => {
                console.log('-error-', error);
                window.location.href = MyConfig.BASE_PATH + '#/login';
                return;
            });
        //findleaf(state.data);
    },[]);

    const selectNode = (e,nodeid) => { 
        //console.log('menu select: ' + nodeid);
        setDom({
            select: {
                id: nodeid, 
                name: state.lobj[nodeid].name,
            },
            evaluation: (dom.select.id === nodeid)?dom.evaluation:null,
            mask: dom.mask,
        });
    };

    if (state.data === null ) {
        return <div>...</div>;
    } else {      
        return (
            <>
                <TreeView
                    className={classes.root}
                    defaultCollapseIcon={<ExpandMoreIcon />}
                    defaultExpanded={state.parents}
                    //expanded={[0,1,2,5]}
                    selected={[dom.select.id]}
                    defaultExpandIcon={<ChevronRightIcon />}
                    onNodeSelect={selectNode}
                >
                    {renderTree(state.data)}
                </TreeView>
            </>
        );
    }
};

export default MenuDomaine;
