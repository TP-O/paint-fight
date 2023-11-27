use log::debug;

use super::{BaseEntity, Entity, EntityId, Position, Size};

pub struct Player {
    base: BaseEntity,
}

impl Entity for Player {
    fn base(&self) -> &BaseEntity {
        &self.base
    }
}

impl Player {
    pub fn new(id: EntityId, size: Size, position: Position) -> Self {
        Self {
            base: BaseEntity { id, size, position },
        }
    }

    pub fn attack(&self) {
        debug!("Player {} attacks", self.base.id())
    }

    pub fn moves(&mut self, position: Position) {
        self.base.set_position(position);

        debug!("Player {} attacks", self.base.id())
    }
}
