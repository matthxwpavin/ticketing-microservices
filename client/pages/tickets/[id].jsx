import useRequest from '../../hooks/use-request';
import Router from 'next/router';

const TicketShow = ({ currentUser, ticket }) => {
    const { doRequest, errors } = useRequest({
        url: '/api/orders',
        method: 'post',
        body: { ticketId: ticket.id },
        onSuccess: order => {
            Router.push(`/orders/${order.orderID}`)
        }
    });
    return (
        <div>
            <h1>{ticket.title}</h1>
            <h4>Price: {ticket.price / 100}</h4>
            <button onClick={_ => doRequest()} className='btn btn-primary'>Purchase</button>
        </div>
    );
}

TicketShow.getInitialProps = async (context, client, currentUser) => {
    const { id } = context.query;
    const { data } = await client.get(`/api/tickets/${id}`);
    return { ticket: data };
};

export default TicketShow;

