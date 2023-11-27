use crate::entity::Position;

pub enum ActionType {
    Attack,
    Move(Position),
}
