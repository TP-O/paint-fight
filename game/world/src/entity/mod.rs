pub mod player;

pub type EntityId = String;
pub type Position = (u16, u16);
pub type Size = (u16, u16);

pub trait Entity {
    fn base(&self) -> &BaseEntity;
}

pub struct BaseEntity {
    id: EntityId,
    size: Size,
    position: Position,
}

impl BaseEntity {
    pub fn id(&self) -> &EntityId {
        &self.id
    }

    pub fn position(&self) -> Position {
        self.position
    }

    fn set_position(&mut self, position: Position) {
        self.position = position;
    }

    pub fn size(&self) -> Size {
        self.size
    }

    fn set_size(&mut self, size: Size) {
        self.size = size;
    }
}
