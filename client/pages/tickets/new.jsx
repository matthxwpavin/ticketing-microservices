import { useState } from 'react';
import Router from 'next/router';
import useRequest from '../../hooks/use-request';

export default () => {
    const [title, setTitle] = useState('');
    const [price, setPrice] = useState('');
    const { doRequest, errors } = useRequest({
        body: {
            title,
            price: Math.floor(parseFloat(price) * 100),
        },
        method: 'post',
        url: '/api/tickets',
        onSuccess: _ => Router.push('/')
    });

    const formatPrice = _ => {
        const value = parseFloat(price).toFixed(2);
        if (Number.isNaN(value)) {
            return;
        }
        setPrice(value);
    }

    const onSubmit = e => {
        e.preventDefault();
        doRequest();
    }

    return (
        <div>
            <h1>Create a ticket</h1>
            <form onSubmit={onSubmit}>
                <div className="form-group mb-3">
                    <label className="form-label">Title</label>
                    <input className="form-control" value={title} onInput={e => setTitle(e.target.value)} />
                </div>
                <div className="form-group mb-3">
                    <label className="form-label">Price</label>
                    <input
                        className="form-control"
                        value={price}
                        onInput={e => setPrice(e.target.value)}
                        type='number'
                        onBlur={formatPrice}
                        inputMode='decimal'
                    />
                </div>
                <button className="btn btn-primary" >Submit</button>
            </form>
        </div>
    );
}