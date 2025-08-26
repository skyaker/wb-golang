import streamlit as st
import requests

API_BASE = "http://order_service:8080"

st.set_page_config(page_title="Order Info", layout="wide")
st.title("Order Info")

order_uid = st.text_input("Enter order_uid")

if st.button("Get order info"):
    if not order_uid:
        st.warning("Enter order_uid")
    else:
        try:
            res = requests.get(f"{API_BASE}/orders/{order_uid}")
            if res.status_code == 200:
                data = res.json()

                tab1, tab2, tab3, tab4 = st.tabs(["Order", "Delivery", "Payment", "Items"])

                with tab1:
                    order = data["order"]
                    st.markdown("### Order")
                    st.markdown(
                        f"""
                        <p style='font-size:18px'>
                        <b>UID:</b> {order["order_uid"]}<br>
                        <b>Track number:</b> {order["track_number"]}<br>
                        <b>Entry:</b> {order["entry"]}<br>
                        <b>Locale:</b> {order["locale"]}<br>
                        <b>Customer ID:</b> {order["customer_id"]}<br>
                        <b>Shard key:</b> {order["shardkey"]}<br>
                        <b>SM ID:</b> {order["sm_id"]}<br>
                        <b>Date created:</b> {order["date_created"]}<br>
                        <b>OOF shard:</b> {order["oof_shard"]}
                        </p>
                        """,
                        unsafe_allow_html=True
                    )

                with tab2:
                    delivery = data["delivery"]
                    st.subheader("Adress")

                    address_parts = []
                    if delivery.get("city"):
                        address_parts.append(delivery["city"])
                    if delivery.get("address"):
                        address_parts.append(delivery["address"])
                    if address_parts:
                        st.text(", ".join(address_parts))

                    st.text(f"Name: {delivery['name']}")
                    if delivery.get("phone"):
                        st.text(f"Phone: {delivery['phone']}")
                    if delivery.get("email"):
                        st.text(f"E-mail: {delivery['email']}")

                with tab3:
                    payment = data["payment"]
                    st.markdown("### Payment")
                    st.markdown(
                        f"""
                        <p style='font-size:18px'>
                        <b>Transaction:</b> {payment["transaction"]}<br>
                        <b>Request ID:</b> {payment["request_id"]}<br>
                        <b>Currency:</b> {payment["currency"]}<br>
                        <b>Provider:</b> {payment["provider"]}<br>
                        <b>Amount:</b> {payment["amount"]}<br>
                        <b>Payment DT:</b> {payment["payment_dt"]}<br>
                        <b>Bank:</b> {payment["bank"]}<br>
                        <b>Delivery cost:</b> {payment["delivery_cost"]}<br>
                        <b>Goods total:</b> {payment["goods_total"]}<br>
                        <b>Custom fee:</b> {payment["custom_fee"]}
                        </p>
                        """,
                        unsafe_allow_html=True
                    )

                with tab4:
                    items = data["items"]
                    st.subheader("Items")
                    if items:
                        st.table([
                            {
                                "ID": it["chrt_id"],
                                "Name": it["name"],
                                "Price": it["price"],
                                "Sale": it["sale"],
                                "Brand": it["brand"],
                                "Status": it["status"],
                                "Track": it["track_number"]
                            }
                            for it in items
                        ])
                    else:
                        st.info("No items")

            elif res.status_code == 404:
                st.error("Order not found")
            else:
                st.error(f"Error: {res.status_code}")
        except Exception as e:
            st.error(f"Error: {e}")
