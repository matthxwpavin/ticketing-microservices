const OrderList = ({ orders }) => {
    const OrderLi = () => orders.map(order => (<li key={order.ID}>
        {order.ticket.title} - {order.status}
    </li>));

    return <ul>
        <OrderLi />
    </ul>;
}

OrderList.getInitialProps = async (context, client) => {
    const { data } = await client.get('/api/orders');
    return { orders: data };
}

export default OrderList;