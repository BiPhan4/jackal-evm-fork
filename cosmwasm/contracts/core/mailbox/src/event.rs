use cosmwasm_std::{Addr, Event, HexBinary};
use hpl_interface::types::Message;

pub fn emit_instantiated(owner: Addr) -> Event {
    Event::new("mailbox_instantiated").add_attribute("owner", owner)
}

pub fn emit_default_ism_set(owner: Addr, new_default_ism: Addr) -> Event {
    Event::new("mailbox_default_ism_set")
        .add_attribute("owner", owner)
        .add_attribute("new_default_ism", new_default_ism)
}

pub fn emit_process_id(id: HexBinary) -> Event {
    Event::new("mailbox_process_id").add_attribute("message_id", id.to_hex())
}

pub fn emit_process(origin: u32, sender: HexBinary, recipient: HexBinary) -> Event {
    Event::new("mailbox_process")
        .add_attribute("origin", format!("{origin}"))
        .add_attribute("sender", sender.to_hex())
        .add_attribute("recipient", recipient.to_hex())
}
