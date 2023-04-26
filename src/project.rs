use std::path::Path;

#[derive(Clone, Debug)]
pub struct Project {
    pub name: String,
    pub configurations: Vec<Configuration>,
}

impl Project {
    pub fn new(name: String, configurations: Vec<Configuration>) -> Project {
        Project {
            name,
            configurations,
        }
    }
}

#[derive(Clone, Debug)]
pub enum Identifier {
    File(String),
    FileExtensionInDirectory(String, String),
}

#[derive(Clone, Debug)]
pub struct Configuration {
    pub identifier: Identifier,
    pub relative_targets: Vec<String>,
}

impl Configuration {
    pub fn new(identifier: Identifier, relative_targets: Vec<String>) -> Configuration {
        Configuration {
            identifier,
            relative_targets,
        }
    }

    pub fn matches(&self, path: &Path) -> bool {
        match &self.identifier {
            Identifier::File(f) => {
                let file_name = path.file_name();
                if let Some(file_name) = file_name {
                    let file_name = file_name.to_str();
                    if let Some(file_name) = file_name {
                        f == file_name
                    } else {
                        false
                    }
                } else {
                    false
                }
            }
            Identifier::FileExtensionInDirectory(dir, ext) => {
                let extension = path.extension();
                let dir_to_compare_to = path.parent();
                match (extension, dir_to_compare_to) {
                    (Some(extension), Some(dir_to_compare_to)) => {
                        let extension = extension.to_str();
                        if let Some(extension) = extension {
                            ext == extension && dir_to_compare_to.ends_with(dir)
                        } else {
                            false
                        }
                    }
                    _ => false,
                }
            }
        }
    }
}
